using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using Newtonsoft.Json;

namespace jsonTest1
{
    class IamHereAns
    {
        public string name;

        [JsonProperty(PropertyName = "JSONRPCPort")]
        public int jsonrpcPort;

        private int value = 1234;
    }

    class Program
    {
        static void Main(string[] args)
        {
            var iamhere = new IamHereAns() { jsonrpcPort = 4444, name = "androidpush.pq" };
            var jsonstr = JsonConvert.SerializeObject(iamhere);
            Console.WriteLine(jsonstr);
        }
    }
}
