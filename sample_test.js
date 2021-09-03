import tcptl from 'k6/x/tlstcp';
import encoding from 'k6/encoding';
import { check } from 'k6';

const server_chain = open('./server_chain.pem')

const tcp_request = "IgQIABAy";
const tcp_request_proto = encoding.b64decode(tcp_request);

export let options = {
  tlsAuth:[
    {
      domains:[""],
      cert: open('./client.pem'),
      key: open('./client_key.pem')

    }
  ],
  tlsVersion: 'tls1.3',
  insecureSkipTLSVerify: true,
  tlsCipherSuites: [
    'TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384',
    'TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256',
  ],
};

const first_url = "20.81.92.115:8057";
export default function () {
  const socket = tcptl.connect("tcp",first_url,server_chain)
  //console.log(socket);
  //TODO:
  //len(tcp_request_proto)
  //create 4 byte array and cast size to it
  //append to array and send


  // check(null, {
  //     'is encoding correct': () => encoding.b64encode(str) === encoding});
    tcptl.send(socket,tcp_request_proto)
    //console.log("tls send passed")
    var eh_receive = tcptl.receive(socket,tcp_request_proto)
    // if(eh_receive > 0)
    // {
    //    console.log("received"+receive_status)
    // }else{
    //    console.log("didn't receive");
    // }
    tcptl.closeConn(socket)
}

export function teardown () {
}
