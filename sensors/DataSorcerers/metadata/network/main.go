package main;

import (
    "fmt"
    "time"
    "context"
    "slices"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"

    "google.golang.org/grpc" 
    "github.com/golang/protobuf/ptypes"

    fbpacket "FusionBridge/metadata/packet"
    bridge "FusionBridge/metadata/network"

    "DataSorcerers/helpers"
)

const (
    snapLen = 1024
)

var NetworkSliceLength uint;
var opts []grpc.DialOption;
var conn *grpc.ClientConn;
var client bridge.NetworkClient;
var err error; 
var NetworkSlice []fbpacket.Packet;
var NetworkSliceIndex uint = 0;


func main() {

    helpers.SetFileName("Network")

    // Load config options
    helpers.LoadConfig()
    if !helpers.Config.RunNetwork { return }


    NetworkSliceLength = helpers.Config.NetworkSliceLength

    NetworkSlice = make([]fbpacket.Packet, NetworkSliceLength)

    conn = helpers.CreateConnection()
    client = bridge.NewNetworkClient(conn); 


     // This will find more devs than in ID. Feature, not bug
    interfac, err := pcap.FindAllDevs()
    if err != nil {
        panic(err)
    }
    for _, el := range interfac {
        inact_handle, err := inactiveSetUp(el.Name)
        if err != nil {
            out := fmt.Sprintf("Error creating inactive handle on dev %v: %v",
                        el.Name, err)
            helpers.LogInfo(out)
            continue
        }
        defer inact_handle.CleanUp()
        handle, err := inact_handle.Activate()
        if err != nil {
            out := fmt.Sprintf("Error activating handle on dev %v: %v",
                    el.Name, err)
            helpers.LogInfo(out)
            continue
        }
        defer handle.Close()
        source := gopacket.NewPacketSource(handle, handle.LinkType())
        go read_packets(el.Name, source)
    }
}


func inactiveSetUp(dev string) (*pcap.InactiveHandle, error) {
    inact_handle, err := pcap.NewInactiveHandle(dev);
    if err != nil {
        out := fmt.Sprintf("Error setting up new inactive handle: %v", err)
        helpers.FatalError(out)
        return nil, err
    }

    // Can strip errors with: inact_handle.Error()
    inact_handle.SetBufferSize(150000); // in bytes

    // packets are delivered to application directly ASAP
    // overrides SetTimeout
    inact_handle.SetImmediateMode(true); 

    // inact_handle.SetPromisc(true);

    // Same idea as SetPromisc but for wireless networks
    // inact_handle.SetRFMon(true);

    // Set read timeout for the handle
    // inact_handle.SetTimeout(1000)

    // Tell pcap how to set timestamps. Idk what these are though
    // ts, err := pcap.TimestampSourceFromString("i have no idea")
    // inact_handle.SetTimestampSource(ts)

    // List supported timestamp sources:
    ts_src_slice := inact_handle.SupportedTimestamps()
    

    adapt_src, err := pcap.TimestampSourceFromString("adapter")
    if err != nil { 
        out := fmt.Sprintf("Error timestamp source from string: %v", err)
        helpers.FatalError(out)
    }

    host_src, err := pcap.TimestampSourceFromString("host")
    if err != nil { 
        out := fmt.Sprintf("Error timestamp source from string: %v", err)
        helpers.FatalError(out)
    }

    if slices.Contains(ts_src_slice, adapt_src) {
        inact_handle.SetTimestampSource(adapt_src)
    } else if slices.Contains(ts_src_slice, host_src) {
        inact_handle.SetTimestampSource(host_src)
    }

    return inact_handle, nil
}


func read_packets(name string, source *gopacket.PacketSource) {
    for pkt := range source.Packets() {
        if err := pkt.ErrorLayer(); err != nil {
            out := fmt.Sprintf("Error reading packet on %v: %v", name, err)
            helpers.LogInfo(out)
            continue
        }

        timestamp := pkt.Metadata().CaptureInfo.Timestamp

        packet_to_send := fbpacket.Packet{
            ApplicationLayer: fbpacket.ApplicationLayer{
                TransportLayer: fbpacket.TransportLayer{
                    NetworkLayer: fbpacket.NetworkLayer{
                        LinkLayer: fbpacket.LinkLayer{
                            PhysicalLayer: fbpacket.PhysicalLayer{
                                TimeStamp : timestamp,
                                Dev       : name,
                            },
                        },
                    },
                },
            },
        };

        if arpPkt := pkt.Layer(layers.LayerTypeARP); arpPkt != nil {
            packet_to_send.LinkLayer.LinkProtocol = "ARP"

            arp := arpPkt.(*layers.ARP)

            packet_to_send.LinkLayer.ARP = fbpacket.ARPPacket{
                 SrcHwAddy  : arp.SourceHwAddress,
                 SrcProtAdd : arp.SourceProtAddress,
                 DstHwAddy  : arp.DstHwAddress,
                 DstProtAdd : arp.DstProtAddress,
                 Operation  : arp.Operation,
                 Protocol   : arp.Protocol,
            }
        } else if ethPkt := pkt.Layer(layers.LayerTypeEthernet); ethPkt != nil {
            packet_to_send.LinkLayer.LinkProtocol = "Ethernet"

            eth := ethPkt.(*layers.Ethernet)

            packet_to_send.LinkLayer.ETH = fbpacket.ETHPacket{
                SRC_MAC : eth.SrcMAC.String(),
                DST_MAC : eth.DstMAC.String(),
                Length  : eth.Length,
            }
        }

        if ipv4 := pkt.Layer(layers.LayerTypeIPv4); ipv4 != nil {
            packet_to_send.NetworkLayer.NetworkProtocol = "IPv4"

            ip := ipv4.(*layers.IPv4)
            // net.IP

            packet_to_send.NetworkLayer.IP = fbpacket.IPPacket{
                V4     : true,
                SRC_IP : ip.SrcIP.String(),
                DST_IP : ip.DstIP.String(),
                // IPProtocol. uint8
                // Protocol := ip.Protocol
            }
        } else if ipv6 := pkt.Layer(layers.LayerTypeIPv6); ipv6 != nil {
            packet_to_send.NetworkLayer.NetworkProtocol = "IPv6"

            ip := ipv6.(*layers.IPv6)

            packet_to_send.NetworkLayer.IP = fbpacket.IPPacket{
                V4     : false,
                SRC_IP : ip.SrcIP.String(),
                DST_IP : ip.DstIP.String(),
            }
        }

        if tcp := pkt.Layer(layers.LayerTypeTCP); tcp != nil {
            packet_to_send.TransportLayer.TransportProtocol = "TCP"

            pkt := tcp.(*layers.TCP)

            packet_to_send.TransportLayer.TCP = fbpacket.TCPPacket{
                SrcPort : pkt.SrcPort,
                DstPort : pkt.DstPort,
            }

        } else if udp := pkt.Layer(layers.LayerTypeUDP); udp != nil {
            packet_to_send.TransportLayer.TransportProtocol = "UDP"

            pkt := udp.(*layers.UDP)

            packet_to_send.TransportLayer.UDP = fbpacket.UDPPacket{
                SrcPort : pkt.SrcPort,
                DstPort : pkt.DstPort,
            }
        }

        if icmpPkt := pkt.Layer(layers.LayerTypeICMPv4); icmpPkt != nil {
            packet_to_send.TransportLayer.TransportProtocol = "ICMP"

            // icmp := icmpPkt.(*layers.ICMPv4)

            packet_to_send.TransportLayer.ICMP = fbpacket.ICMPPacket{}

        } else if icmp := pkt.Layer(layers.LayerTypeICMPv6); icmp != nil {
            packet_to_send.TransportLayer.TransportProtocol = "ICMP"

            // icmp := icmpPkt.(*layers.ICMPv4)

            packet_to_send.TransportLayer.ICMP = fbpacket.ICMPPacket{}
        }

        if tls := pkt.Layer(layers.LayerTypeTLS); tls != nil {
            packet_to_send.ApplicationProtocol = "TLS"
            packet_to_send.TLS = fbpacket.TLSPacket{}
        }


        if dnsPkt := pkt.Layer(layers.LayerTypeDNS); dnsPkt != nil {
            packet_to_send.ApplicationProtocol = "DNS"

            dns := dnsPkt.(*layers.DNS)

            // Minimal info from the DNS queries
            packet_to_send.DNS = fbpacket.DNSPacket{
                Questions : dns.Questions,
                Answers   : dns.Answers,
            }
        }
        add_to_network_slice(packet_to_send)
    }
}


func add_to_network_slice(pkt fbpacket.Packet) {

    NetworkSlice[NetworkSliceIndex] = pkt;

    NetworkSliceIndex = (NetworkSliceIndex + 1) % NetworkSliceLength;
    if (NetworkSliceIndex == 0) { save_network_slice(); }
}


func save_network_slice() {

    ctx, cancel := context.WithTimeout(context.Background(), 
        10*time.Second)
    defer cancel()

    timestamp, _ := ptypes.TimestampProto(time.Now())

    BridgeNetworkSlice := make([]*bridge.Packet, NetworkSliceLength, NetworkSliceLength);

    for i, pkt := range NetworkSlice {
        BridgeNetworkSlice[i] = fbpacket.ToBridgePacket(pkt)
    }

    to_send  := bridge.NetworkSlice{
        SubmissionNumber:0,
        TimeStamp: timestamp,
        Packets: BridgeNetworkSlice,         
        DevID: helpers.GetID(),
    }

    _, e := client.LogNetworkActivity(ctx, &to_send)
    if e != nil {
        out := fmt.Sprintf("client error in client.metadata.network.send_data: %v", e)
        helpers.LogInfo(out)
    }
}

