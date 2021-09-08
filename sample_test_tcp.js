import tcptls from 'k6/x/tlstcp';
import { check } from 'k6';
import { Rate } from 'k6/metrics';
const first_url = "20.81.92.115:8057";

const server_chain = open('./server_chain.pem')

const tcp_request = "IgQIABAy";
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
  thresholds: {
    errors: ['rate<0.1'], // <10% errors
  },
};
export let errorRate = new Rate('errors');


// 1. init code


// 2. setup code
export function setup()
{
  
}
// 3. VU code

export default function()
{
	const conn =  tcptls.connect("tcpv4",first_url,server_chain)
  
  tcptls.send(conn,"GwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARwAAJAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
  const eh_receive = tcptls.receive(conn)
  // const socket = tcptl.connect("tcp",first_url,server_chain)
  // tcptl.send(socket,tcp_request)
  // //console.log("tls send passed")
  // var eh_receive = tcptl.receive(socket)
  console.log(eh_receive.length)
  
  const res = check(eh_receive, {
    'received bytes': (r) => r.length > 0,
  });
  
  errorRate.add(!res);
    // //todo decode protobuffer
    
  tcptls.closeConn(socket)
}
// 4. teardown code
export function teardown ()
{
}
