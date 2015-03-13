// marcoPolo project
// Copyright 2013 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
'marco polo': 
my own mini zero-configuration networking (just 1-2 'hardcoded' udp port(s) ;-p). 


v0.2 with marcoPolo server (simpler)
--------------------------
A 'marcoPolo server' listens on a (broadcast) UDP port for messages.
It serves as a central point for apps to find each other.
** Just 'marco polo', : app registers itself, others an query for app address

+ possible broadcast UDP msgs ?

see mindmap marcoPoloSimpler.mm
(freemind)


v0.2 with marcoPolo server (scrapped)
--------------------------
A 'marcoPolo server' listens on a (broadcast) UDP port for messages.
It serves as a central point for apps to find each other.

An app (polo) registers messages it will wants to receive, ie 'services' it offers,
other apps (marco) requests marcoPolo for the 'service'/msg, 
marcoPolo answers with the address+port of the registered app.
The two apps then communicate with each other w/o any further marcoPolo calls.

It is also possible to send broadcast messages.
An app sends a msg to marcoPolo which will then fwd it to all apps that registered for it

see mindmap marcoPoloSimple.mm
(freemind)


v0.1 direct PC-marco to android-polo
------------------------------------
PC calls 'marco', android answers 'polo' on UDP, 
this way the PC can find the Android IP.

Android/Polo waits for udp datagram 'marco' on port 4444
PC/Polo broadcasts 'marco' on UDP to port 4444 (periodically, in a loop)
polo recvs 'marco' and answers 'polo' to PC using the address found in te datagram (ie we know who called us)
marco recvs 'polo' from Android, thus we know the IP address of the Android device.

Currently
marco sends "marco|marcoPoloId"
polo responds "polo|marcoPoloId|poloListeningPort"

which means polo acts as a server once the marcoPolo 'protocol' is completed, 
marco then needs to connect to polo.

#TODO: 
+ change marco msg to also include a server port (on which marco receives connections after marcoPolo protocol)
 => "marco|marcoPoloId|marcoListeningPort"
this way it's possible to have either marco or polo act as a server (does both make sense !?)
(up to app writer to decide how to use this)

Note that normally / currenlty the 'server' is on the polo side. 
so to have, for example, polo be on the PC side, we need to have marco Android code...
As of now I have a C++ lib that handles both marco & polo.
But in Android JAVA I just have a polo example (polo is much simpler to implement)

The C++ lib uses boost for asynch socket I/O
