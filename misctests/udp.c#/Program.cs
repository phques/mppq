using System;
using System.Net;
using System.Net.Sockets;
using System.Text;
namespace udpbroadcast
{
    class Program
    {
        static void Main(string[] args)
        {
            Socket s = new Socket(AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);
            s.EnableBroadcast = true;

            //IPAddress broadcast = IPAddress.Parse("255.255.255.255");
            IPAddress broadcast = IPAddress.Parse("192.168.1.255");
            IPEndPoint ep = new IPEndPoint(broadcast, 4444);

            //var udpClient = new UdpClient(ep);

            //byte[] sendbuf = Encoding.ASCII.GetBytes(args[0]);
            byte[] sendbuf = Encoding.ASCII.GetBytes("quit");
            s.SendTo(sendbuf, ep);

            Console.WriteLine("Message sent to the broadcast address");

            Console.WriteLine(s.LocalEndPoint);   // 0.0.0.0
        }
    }
}
