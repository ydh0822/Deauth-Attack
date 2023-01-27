# Deauth-Attack

BoB 11th - Digital Forensics Track
H4uN

Usage : 
You must run with sudo privileges to obtain iwconfig privileges.

Deauth-Attack 3ways
 - AP broadcast frame
    - ./Deauth-Attack <interface> <ap mac>
 - AP unicast, Station unicast frame
    - ./Deauth-Attack <interface> <ap mac> <station mac>
 - authentication frame
    - ./Deauth-Attack <interface> <ap mac> <station mac> -c