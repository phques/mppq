// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
//--marcoPolo.cpp --
#include "stdafx.h"
#include <boost/algorithm/string/classification.hpp>
#include <boost/thread.hpp>
#include <boost/bind.hpp>
#include <boost/lexical_cast.hpp>
#include "marcoPolo.h"

//#include <boost/algorithm/string/classification.hpp>

using namespace boost::asio;
using namespace boost::algorithm;
using namespace std;

using boost::asio::ip::udp;

const char* MSG_SEPAR = "|";


MarcoPolo::MarcoPolo(io_service& ioService, string poloName, unsigned short poloPort)
    : ioService(ioService), poloName(poloName), poloPort(poloPort), foundPolo(false)
{
}

MarcoPolo::~MarcoPolo()
{
}


//--------------

// Send a 'marco' message (broadcast udp) looking for a 'polo' answer
// this lets us find a remote app/pier with minimum config/parameters
bool MarcoPolo::marco()
{
    // loop on whole thing (re-send & recv, if not right answer from polo)
    foundPolo = false;

    // setup local socket used to send data, broadcast (any port)
    udp::socket sock(ioService, udp::endpoint(udp::v4(), 0));
    socket_base::broadcast option(true);
    sock.set_option(option);

    while (!foundPolo)
    {
        //-- launch asynch recv polo response --

        sock.async_receive_from(
            boost::asio::buffer(data, max_length), responseEndpoint,
            boost::bind(&MarcoPolo::handleReceiveFromPolo, this,
                   boost::asio::placeholders::error,
                   boost::asio::placeholders::bytes_transferred));

        bool gotData = false;
        while (!foundPolo && !gotData)
        {
            // send a(nother) 'marco'
            sendMarco(sock);

            // wait up to 1sec, checking if we recvd data
            for (int i = 0; i < 10 && !gotData; i++) {
                // check for recvd data after waiting a bit
                boost::this_thread::sleep(boost::posix_time::milliseconds(100));

                if (ioService.poll_one() >= 1)
                    gotData = true;
            }
        }
    }

    return foundPolo;
}


void MarcoPolo::sendMarco(udp::socket& sock)
{
    //--- send 'marco', broadcast UDP --

    udp::endpoint remoteEndpt(ip::address_v4::broadcast(), poloPort);

    string marcoMsg = this->marcoMsg();

    //##debug
    cout << "Broadcasting " << marcoMsg << " on " << sock.local_endpoint() << " to port " << poloPort << "\n";

    sock.send_to(boost::asio::buffer(marcoMsg), remoteEndpt);

}

// handle asynch recv, recv 'polo' answer for our 'marco' call
void MarcoPolo::handleReceiveFromPolo(const boost::system::error_code& error, size_t bytes_recvd)
{
    if (!error && bytes_recvd > 0)
    {
        // answr was read int 'data'
        data[bytes_recvd] = 0;
        string reply(data);

        //##debug
        cout << "Reply is: " << reply << "\n";
        cout << "from " << responseEndpoint << endl;

        // decode reply !
        // split by separator MSG_SEPAR
        vector<string> responseParts;
        split(responseParts, reply, is_any_of(MSG_SEPAR), token_compress_on);
        if (responseParts.size() == 3)
        {
            if (responseParts[0] == "polo" && responseParts[1] == poloName)
            {
                //##TODO: check that all digits string / convert to unsigned short
                using boost::lexical_cast;
                using boost::bad_lexical_cast;
                try
                {
                    // try to convert port to unsigned short
                    poloResponsePort_ = (lexical_cast<unsigned short >(responseParts[2]));
                    foundPolo = true;
                } catch(bad_lexical_cast &) {
                    //## debug
                    cout << "not a valid port recvd from polo : " << responseParts[2] << endl;
                }
            }
        }
    }

    if (!foundPolo)
    {
        //## debug
        cout << "got error or empty / non-valid polo reply" << endl;
    }
}


//--------------

// wait for a 'marco' message asking for our 'polo'
// send back 'polo' with our tcp listen socket port
// leave it to the caller to loop & call us again if we exit with false
bool MarcoPolo::polo(unsigned short poloListenTcpPort)
{
    // -- recv --

    // Receive UDP datagram from socket sock, return as a string
    // callEndpoint will hold the endpoint of the 'caller'
    udp::endpoint callEndpoint;
    udp::socket sock(ioService, udp::endpoint(udp::v4(), poloPort));

    char data[max_length+1];
    size_t length = sock.receive_from(
        boost::asio::buffer(data, max_length), callEndpoint);

    data[length] = 0;
    string msg(data);

    //## debug
    cout << "received  : " << msg << "\n";
    cout << "from " << callEndpoint << endl;

    // check if a marco asking for us!
    bool ok = (msg == marcoMsg());
    if (ok) {
        // -- send back 'polo' answer to 'marco' --
        string poloMsg = this->poloMsg(poloListenTcpPort);
        sock.send_to(boost::asio::buffer(poloMsg), callEndpoint);
    }
    else {
        //## debug
        cout << "NOT a MARCO request for us" << endl;
    }

    return ok;
}

//--------------

// the 'marco' text message to send
string MarcoPolo::marcoMsg()
{
    stringstream msg;
    msg << "marco" << MSG_SEPAR << poloName;

    return msg.str();
}

// the 'polo' reply text message to send
string MarcoPolo::poloMsg(unsigned short poloTcpPort)
{
    stringstream msg;
    msg << "polo" << MSG_SEPAR << poloName << MSG_SEPAR << poloTcpPort;

    return msg.str();
}
