// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
// -- mainMarco.cpp --

#include "stdafx.h"
#include "marcoPolo.h"

using namespace std;
using namespace boost::asio;

using boost::asio::ip::udp;



int main(int argc, char* argv[])
{
    try
    {
        //    if (argc != 3)
        //    {
        //      cerr << "Usage: blocking_udp_echo_client <host> <port>\n";
        //      return 1;
        //    }

        io_service io_service;

        MarcoPolo marcoPolo(io_service, "testMarcoPolo");

        cout << marcoPolo.marco() << "\n";
        cout << "found app port is " << marcoPolo.poloResponsePort() << endl;

        boost::asio::ip::udp::endpoint poloEndpoint = marcoPolo.poloEndpoint();
        cout << "app is on " << poloEndpoint.address() << endl;
    }
    catch (exception& e)
    {
        cerr << "Exception: " << e.what() << "\n";
    }

    return 0;
}
