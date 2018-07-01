//=====================================
// Unreal 2 Masterserver
// Proof of concept.
//
// Half of this code may work, 
// half of it may crash the game.
//
// At least, it works for Pariah/Warpath...
//
// ONLY SUPPORT 1 CLIENT AT A TIME.
// WHENEVER ANOTHER CLIENT CONNECTS, IT BREAKS THE OLD ONE!
//=====================================

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;

namespace ConsoleApp1
{

    //======================
    // Pariah
    // Warpath
    //======================
    class Pariah_MasterServer : Unreal_TCPListener
    {
        public override char ClientCharacter { get { return System.Convert.ToChar(System.Convert.ToUInt32("005C", 16)); } }
        public override char ServerCharacter { get { return System.Convert.ToChar(System.Convert.ToUInt32("0060", 16)); } }
        public override string ClientName { get { return "PARIAHCLIENT"; } } 
        public override string ServerName { get { return "PARIAHSERVER"; } }
    }

    //======================
    // Land of the dead
    // Day of the Zombie
    //======================
    class LOTD_MasterServer : Unreal_TCPListener
    {
        public override char ClientCharacter { get { return 'V'; } }
        public override char ServerCharacter { get { return '<'; } }
        public override string ClientName { get { return "CLIENT"; } }
        public override string ServerName { get { return "SERVER"; } }
    }

    class Unreal_TCPListener
    {

        public virtual char ClientCharacter { get; set; }
        public virtual char ServerCharacter { get; set; }
        public virtual string ClientName { get; set; }
        public virtual string ServerName { get; set; }

        // When the server sends this,
        // the client should return something like this:
        //  V   !<hash32> !<hash32> CLIENT  int
        private byte[] MSG_AUTHENTIFICATION = { 0x03, 0x00, 0x00, 0x00, 0x02, 0x30, 0x00 };         // authentification request
        
        // Those bytes are sent as a response to the authentification.
        private byte[] MSG_DENIED = { 0x08, 0x00, 0x00, 0x00, 0x09, 0x44, 0x45, 0x4e, 0x49, 0x45, 0x44, 0x00 };
        private byte[] MSG_VERIFIED = { 0x0a, 0x00, 0x00, 0x00, 0x09, 0x41, 0x50, 0x50, 0x52, 0x4f, 0x56, 0x45, 0x44, 0x00 };

        //private byte[] MSG_GAMETYPE_REQUEST = { 0x0e, 0x00, 0x00, 0x00, 0x00, 0x01, 0x09, 0x67, 0x61, 0x6d, 0x65, 0x74, 0x79, 0x70, 0x65, 0x00, 0x00, 0x01 }; // ? was probably sent on LOTD/DOTZ
        private byte[] MSG_GAMETYPE_REQUEST = { 0x01, 0x00, 0x00, 0x00, 0x01 };
        private byte[] MSG_GAMETYPE_RESPONSE = { 0x05, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x01 };


        // CURRENTLY UNDER INVESTIGATION
        // It was sent by a server to the Masterserver, what is its meaning?
        private byte[] MSG_SERVER_REQUEST = { 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }; // 12 bytes starting with 0x08 ?

        // MOTD request by the client.
        byte[] MSG_MOTD_REQUEST = { 0x01, 0x00, 0x00, 0x00, 0x01 };

        // PAST THIS STATE, the client may request a GAMETYPE message
        // Request to ask the GameType
        // 0e = Probably the header
        // 00 x 4 = Probably to separate, or just the command is using 5 bytes ?
        // "gametype" text as bytes
        // 00 x 2 = ?
        // 01 = Probably to verify the packet checksum is done?

        private bool isHexadecimal(string data)
        {
            bool isHex = data.All("0123456789abcdefABCDEF".Contains);
            return isHex;
        }

        // We're going to check the request status.
        // 0 = DENIED
        // 1 = UPGRADE
        // 2 = CLIENT APPROVED
        // 3 = SERVER APPROVED
        private int CheckVerification(string data)
        {
            //Console.WriteLine(String.Format("[ BASE =====> {2} ] CL ==> {0} \\ SV ==> {1}", ClientCharacter, ServerCharacter, data[0]));

            // Verification Flag
             if (data[0] != ClientCharacter && data[0] != ServerCharacter)
             {
                 Console.WriteLine("Check Code failed?");
                return 0;
             }

             // Verification of Data #1
             if (data[4] != '!')
             {
                 Console.WriteLine("Check ! #1 failed?");
                 return 0;
             }

             if (!isHexadecimal(data.Substring(5, 32)))
             {
                 Console.WriteLine("Check Auth #1 failed?");
                 return 0;
             }

             // ToDo: We won't forget for Data #2

            // Check if we're a client.
            if (data.Contains(ClientName))
                return 1;

            // Check if we're actually a server then ?
            else if (data.Contains(ServerName))
                return 3;

            // No idea about who it is, reject it immediately.
            return 0;
        }

        // Copying bytes code adapted from this page.
        // https://stackoverflow.com/questions/5591329/c-sharp-how-to-add-byte-to-byte-array
        //---------------------------------------------------------------------------------
        public byte[] addByteToArray(byte[] bArray, byte newByte)
        {
            byte[] newArray = new byte[bArray.Length + 1];
            bArray.CopyTo(newArray, 0);
            newArray[bArray.Length] = newByte;
            return newArray;
        }

        public Unreal_TCPListener()
        {
            TcpListener server = null;
            try
            {
                // Set the TcpListener on carPort.
                int port = 27900;
                server = new TcpListener(IPAddress.Any, port);

                // Start listening for client requests.
                server.Start();

                // Buffer for reading data
                Byte[] bytes = new Byte[1024];
                String data = null;

                int iState = 0;

                Console.WriteLine("=======================");
                Console.WriteLine("Unreal Engine 2 Masterserver");
                Console.WriteLine(String.Format("PoC started on port {0}", port));
                Console.WriteLine("=======================");

                // Enter the listening loop.
                while (true)
                {
                    // Perform a blocking call to accept requests.
                    TcpClient client = server.AcceptTcpClient();
                    Console.WriteLine("Looking for a request...");

                    NetworkStream _NetworkStream = client.GetStream();

                    data = null;
                    bool dataAvailable = false;

                    // Get a stream object for reading
                    BinaryReader stream = new BinaryReader(client.GetStream());
                    int receivingBufferSize = (int)client.ReceiveBufferSize;


                    Console.WriteLine(String.Format("<SERVER> Sending Authentification request", port));
                    _NetworkStream.Write(MSG_AUTHENTIFICATION, 0, MSG_AUTHENTIFICATION.Length);

                    int dabytes;
                    while (true)
                    {
                        if (!dataAvailable)
                        {
                            dataAvailable = _NetworkStream.DataAvailable;
                            if (server.Pending())
                            {
                                Console.WriteLine(String.Format("Requester has connected!"));
                                break;
                            }
                        }

                        if (dataAvailable)
                        {
                            // Loop to receive all the data sent by the client.
                            dabytes = stream.Read(bytes, 0, bytes.Length);
                            data = System.Text.Encoding.ASCII.GetString(bytes, 0, dabytes);

                            byte[] BytesList = Encoding.Default.GetBytes(data);
                            var HexToString = BitConverter.ToString(BytesList);

                            Console.WriteLine("<CLIENT> Received \"{0}\" ({1})", data, HexToString);

                            if (iState == 0)
                            {
                                if (CheckVerification(data) == 0)
                                {
                                    _NetworkStream.Write(MSG_DENIED, 0, MSG_DENIED.Length);    // Authentification failed, ba-bai~ ♪
                                    break;
                                }
                                else
                                {
                                    Console.WriteLine("Client has been verified.");
                                    _NetworkStream.Write(MSG_VERIFIED, 0, MSG_VERIFIED.Length);
                                    iState = 1;
                                }
                            }
                            
                            if (iState == 1)
                            {
                                if (BytesList[0] == MSG_MOTD_REQUEST[0] && BytesList[4] == MSG_MOTD_REQUEST[4])   // 1 0 0 0 1 seems to be a MOTD request
                                {
                                    Console.WriteLine("<CLIENT> MOTD Requested");

                                    var byteArray = Encoding.ASCII.GetBytes("Ch0wW, get the f0k to sleep ! You baguette !!!");  // Objective
                                    int iUpdateStatus = 0;

                                    byte[] msganswer = { 0x00 };
                                       
                                    for (int i = 0; i < 3; i++)
                                        msganswer = addByteToArray(msganswer, 0x00);

                                    msganswer = addByteToArray(msganswer, BitConverter.GetBytes(byteArray.Length)[0]);  // Total length of MOTD

                                    foreach (var letter in byteArray)
                                        msganswer = addByteToArray(msganswer, letter);                  // Adding characters as bytes to MOTD

                                    msganswer[0] = BitConverter.GetBytes(msganswer.Length)[0];      // Calculating current MSG length

                                    for (int i = 0; i < 3; i++)
                                        msganswer = addByteToArray(msganswer, 0x00);

                                    if (iUpdateStatus != 0)
                                    {
                                        /*
                                         * I STILL DON'T GET THIS PART so don't expect it to be working
                                         */
                                        msganswer = addByteToArray(msganswer, BitConverter.GetBytes(iUpdateStatus)[0]);
                                        msganswer = addByteToArray(msganswer, BitConverter.GetBytes(4)[0]);

                                        /*foreach (var letter in teststr)
                                            msganswer = addByteToArray(msganswer, letter);                  // Adding characters as bytes to MOTD*/

                                        for (int i = 0; i < 3; i++)
                                            msganswer = addByteToArray(msganswer, 0x00);
                                    }
                                    else
                                    {
                                        msganswer = addByteToArray(msganswer, 0x00);
                                    }
                                    /*  for (int i = 0; i < 3; i++)
                                          msganswer = addByteToArray(msganswer, 0x00);*/

                                    Console.WriteLine(String.Format("<SERVER> Sending \"{0}\" ({1})", System.Text.Encoding.UTF8.GetString(msganswer), BitConverter.ToString(msganswer)));
                                    _NetworkStream.Write(msganswer, 0, msganswer.Length);
                                    // iState = 2;
                                }
                            }

                            dataAvailable = false;
                        }

                        if (server.Pending())
                        {
                            Console.WriteLine("Som");
                            break;
                        }
                    }

                    Console.WriteLine("Client close");
                    // Shutdown and end connection
                    client.Close();
                    iState = 0;         // Restart iState
                }
            }
            catch (SocketException e)
            {
                Console.WriteLine("[CRITICAL] ERROR on SocketException ► {0}", e);
            }
            finally
            {
                // Stop listening for new clients.
                Console.WriteLine("STOPPING SERVER.");
                server.Stop();
            }
        }

    }
}
