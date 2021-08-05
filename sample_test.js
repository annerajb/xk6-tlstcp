import tcptl from 'k6/x/tlstcp';
import encoding from 'k6/encoding';
import { check } from 'k6';
const cert_chain_secp256r1 = `-----BEGIN CERTIFICATE-----

-----END CERTIFICATE-----
`;
const first_server_ed25519_chain = `-----BEGIN CERTIFICATE-----

-----END CERTIFICATE-----

`;
//open('./mycert.pem')

const tcp_request = "EgsKCTAwMDEyMzAwMA==";
const tcp_request_proto = encoding.b64decode(tcp_request);

export let options = {
  tlsVersion: 'tls1.3',
  insecureSkipTLSVerify: true,
  tlsCipherSuites: [
    'TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384',
    'TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256',
  ],
};

const first_url = "asd.eastus.cloudapp.azure.com:8585";
const second_url = "asdasdasd.eastus.cloudapp.azure.com:8589";
const third_crl_url = "asdasdasdasd.eastus.cloudapp.azure.com:8585";
export default function () {
  const socket = tcptl.connect("tcp",first_url,first_server_ed25519_chain)
  //TODO:
  //len(tcp_request_proto)
  //create 4 byte array and cast size to it
  //append to array and send


  // check(null, {
  //     'is encoding correct': () => encoding.b64encode(str) === encoding});
    tcptl.send(socket,gsds_request_proto)
    console.log("already sent")
    var eh_receive = tcptl.receive(socket,gsds_request_proto)
    console.log(eh_receive);
    // if(lenn > 0)
    // {
    //   console.log("received"+receive_status)
    // }else{
    //   console.log("didn't receive");
    // }
    tcptl.closeConn(socket)
}

export function teardown () {
}