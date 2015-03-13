// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
// -- mainPolo.cpp --
#include "stdafx.h"
#include "marcoPolo.h"

using namespace boost::asio;
using namespace std;

using boost::asio::ip::udp;



int main(int argc, char* argv[])
{
    try
    {
//    if (argc != 3)
//    {
//      std::cerr << "Usage: blocking_udp_echo_client <host> <port>\n";
//      return 1;
//    }

        io_service io_service;

        MarcoPolo marcoPolo(io_service, "testMarcoPolo");

        unsigned short poloListenTcpPort = 1234; // would be the port of the socket we do listen() with
        bool ok = false;
        while (true) {
            ok = marcoPolo.polo(poloListenTcpPort);
            cout << ok << endl;
        }
    } catch (exception& e) {
        cerr << "Exception: " << e.what() << "\n";
    }

  return 0;
}
