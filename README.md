# Deauth-Attack

BoB 11th - Digital Forensics Track
H4uN

Usage : 
You must run with sudo privileges to obtain iwconfig privileges.

Deauth-Attack 3ways
 - AP broadcast frame
    - ./Deauth-Attack -i interface -a ap_mac
 - AP unicast, Station unicast frame
    - ./Deauth-Attack -i interface -a ap_mac -s station_mac
 - authentication frame
    - ./Deauth-Attack -i interface -a ap_mac -s station_mac -c