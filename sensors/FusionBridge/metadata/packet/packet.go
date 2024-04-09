package packet;

import (
    "time"
    "github.com/google/gopacket/layers"
    "github.com/golang/protobuf/ptypes"
    "FusionBridge/metadata/network"
)

type PhysicalLayer struct {
    TimeStamp  time.Time // A unix timestamp
    Dev        string    // Where it was captured
}

// I hate this implementation but no inheritence forces my hand
type LinkLayer struct {
    PhysicalLayer
    LinkProtocol string
    ARP ARPPacket
    ETH ETHPacket
}


type NetworkLayer struct {
    LinkLayer
    NetworkProtocol string
    IP IPPacket
    ICMP ICMPPacket
}


type TransportLayer struct {
    NetworkLayer
    TransportProtocol string
    TCP TCPPacket
    UDP UDPPacket
}


type ApplicationLayer struct {
    TransportLayer
    ApplicationProtocol string
    TLS TLSPacket
    DNS DNSPacket
}

type Packet struct {
    ApplicationLayer
}

type ARPPacket struct {
    // Specific to ARP
    SrcHwAddy  []byte
    SrcProtAdd []byte
    DstHwAddy  []byte
    DstProtAdd []byte
    Operation  uint16
    Protocol   layers.EthernetType // unint16
}


type ETHPacket struct {
    // Specific to Eth packets
    SRC_MAC string
    DST_MAC string
    Length  uint16 
}


type IPPacket struct {
    // Specific to IP packets
    V4       bool
    SRC_IP   string
    DST_IP   string
    // SRC_Port int32
    // DST_Port int32
}

type ICMPPacket struct {
    IPPacket // Uses IP for routing so we will compose here
    // Specific to ICMP packets
}

type TCPPacket struct {
    // Specific to TCP packets
    SrcPort layers.TCPPort
    DstPort layers.TCPPort
}


type UDPPacket struct {
    // Specific to UDP packets
    SrcPort layers.UDPPort // uint16 wrap
    DstPort layers.UDPPort
}


type TLSPacket struct {
    // Specific to TLS
}

type DNSPacket struct {
    // Specific to DNS
    Questions []layers.DNSQuestion
    Answers []layers.DNSResourceRecord
}



func ToBridgePacket(in Packet) *network.Packet {
    timestamp, _ := ptypes.TimestampProto(in.PhysicalLayer.TimeStamp)
    // Set up general definitions
    out := network.Packet{
        TimeStamp: timestamp,
        Dev: in.PhysicalLayer.Dev,
        LinkProtocol:in.LinkLayer.LinkProtocol,
        NetworkProtocol:in.NetworkLayer.NetworkProtocol,
        TransportProtocol:in.TransportLayer.TransportProtocol,
        ApplicationProtocol:in.ApplicationLayer.ApplicationProtocol,
    };
    // Fill in corresponding link layer packets
    if in.LinkLayer.LinkProtocol == "ARP" {
        out.ARP = &network.ARPPacket{
            SrcHwAddy   :   in.LinkLayer.ARP.SrcHwAddy,
            SrcProtAdd  :   in.LinkLayer.ARP.SrcProtAdd,
            DstHwAddy   :   in.LinkLayer.ARP.DstHwAddy,
            DstProtAdd  :   in.LinkLayer.ARP.DstProtAdd,
            Operation   :   uint32(in.LinkLayer.ARP.Operation),
            Protocol    :   uint32(in.LinkLayer.ARP.Protocol),
        }
    } else if in.LinkLayer.LinkProtocol == "Ethernet" {
        //out.ETH = in.LinkLayer.ETH
        out.ETH = &network.EthernetPacket{
            SRC_MAC:in.LinkLayer.ETH.SRC_MAC,
            DST_MAC:in.LinkLayer.ETH.DST_MAC,
            Length:uint32(in.LinkLayer.ETH.Length),
        }
    }
    // Fill in corresponding network layer packets
    if in.NetworkLayer.NetworkProtocol == "IP" {
        out.IP = &network.IPPacket{
            V4    : in.NetworkLayer.IP.V4,
            SRC_IP: in.NetworkLayer.IP.SRC_IP,
            DST_IP: in.NetworkLayer.IP.SRC_IP,
        }
    } else if in.NetworkLayer.NetworkProtocol == "ICMP" {
        out.ICMP = &network.ICMPPacket{}
    }
    // Fill in corresponding transport layer packets

    if in.TransportLayer.TransportProtocol == "TCP" {
        out.TCP = &network.TCPPacket{
            SrcPort: uint32(in.TransportLayer.TCP.SrcPort),
            DstPort: uint32(in.TransportLayer.TCP.DstPort),
        }
    } else if in.TransportLayer.TransportProtocol == "UDP" {
        out.UDP = &network.UDPPacket{
            SrcPort: uint32(in.TransportLayer.UDP.SrcPort),
            DstPort: uint32(in.TransportLayer.UDP.DstPort),
        }
    }
    // Fill in corresponding application layer packets
    if in.ApplicationLayer.ApplicationProtocol == "DNS" {
        // Generate the lists of questions & answers
        var questions []*network.DNSQuestion;
        var answers []*network.DNSResourceRecord; 
        for _, el := range in.ApplicationLayer.DNS.Questions {
            questions = append(questions, &network.DNSQuestion{
                Type: uint32(el.Type),
                Name: string(el.Name),
            })
        }
        for _, el := range in.ApplicationLayer.DNS.Answers {
            answers = append(answers, &network.DNSResourceRecord{
                Type: uint32(el.Type),
                Name: string(el.Name),
                IP: string(el.IP.String()),
            })
        }
        // Actually set up those values
        out.DNS = &network.DNSPacket{
            Questions: questions,
            Answers  : answers,
        }
    }
    return &out;
}

