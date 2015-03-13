// AndroidPush project
// Copyright 2012 Philippe Quesnel
// Licensed under the Academic Free License version 3.0
//--marcoPolo.h --
#ifndef MARCO_H
#define MARCO_H

#include <string>
#include <boost/asio.hpp>

class MarcoPolo
{
    enum { max_length = 1024 };

    public:
        MarcoPolo(boost::asio::io_service& ioService, std::string poloName, unsigned short poloPort=4444);
        virtual ~MarcoPolo();

        bool marco();
        bool polo(unsigned short poloListenTcpPort);

        unsigned short poloResponsePort() { return poloResponsePort_; }
        boost::asio::ip::udp::endpoint poloEndpoint() { return responseEndpoint; }

    private:
        std::string recv(boost::asio::ip::udp::socket& sock, boost::asio::ip::udp::endpoint& responseEndpoint);
        void handleReceiveFromPolo(const boost::system::error_code& error, size_t bytes_recvd);
        void sendMarco(boost::asio::ip::udp::socket& sock);

        std::string marcoMsg();
        std::string poloMsg(unsigned short poloListenTcpPort);

    private:
        boost::asio::io_service& ioService;
        std::string poloName;
        unsigned short poloPort;
        bool foundPolo;

        boost::asio::ip::udp::endpoint responseEndpoint;
        unsigned short poloResponsePort_;
        char data[max_length+1];
};


#endif // MARCO_H
