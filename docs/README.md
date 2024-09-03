## problem

Researchers across the US use cyberinfrastructuress to conduct experimental activities. The value
that these infrastructures provide makes them _compelling targets for attackers_, who may wish to
exfiltrate research data, ransom the infrastructure via encryption, or enlist resources into botnets
to send spam, participate in DDOS, or perform cryptomining.

## overview

The DISCERN project at the [USC Information Sciences Institute](https://isi.edu) is producing
datasets to capture the behavior of such malicious activities on research cyberinfrastructure resources.
Our work focuses on [SPHERE testbed](https://sphere-project.net), a novel research infrastructure
which provides resources and services in support of security and privacy research. Our approach is
to conduct malicious activities on the infrastructure in a controlled fashion, alongside legitimate
activities from SPHERE users, and to produce datasets that capture those malicious and benign
activities. Our goal is that these datasets provide visibility into how malicious activities play
out on cyberinfrastructures, thereby enabling research into cyber defenses that can better secure
the nation's cyberinfrastructure resources.

## tools

- [BYOB](https://github.com/STEELISI/byob) - We use BYOB ("Build Your Own Botnet") to construct
  malicious experiments
- [Instrumentation](https://gitlab.com/mergetb/tech/instrumentation") - We are building sensors into
  the Merge testbed platform, which powers the SPHERE research infrastructure. These tools capture
  network metadata, OS metadata, system logs, and user interface activities. Our tools are designed
  for Merge, but are general enough to be deployed on most computational platforms.

## publications and presentations

[2024 CICI PI Meeting](./2024_Brian-Kocoloski_NSF-CICI.pdf) - overview of DISCERN's activities
