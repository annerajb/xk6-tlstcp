import tcptls from 'k6/x/tlstcp';
import { check } from 'k6';
import { Rate } from 'k6/metrics';
const ntp_url = `${__ENV.MY_HOSTNAME}`;//"20.84.29.201:8123";

const server_chain = open('./ntp.pubkey')

export let options = {
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
export default function ()
{
  const conn =  tcptls.connect("udp4",ntp_url,server_chain)

  //ntp protocol sends 7 request then listens for 7 responses
  //then after verifying all the data it averages multiple time differences with the numbers sent and returned
  for(var i = 0; i < 7; ++i)
  {
    const snd = tcptls.send(conn,"GwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARwAAJAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
    //console.log(snd.status);
  }
  
  for(var i = 0; i < 7; ++i)
  {
    const eh_receive = tcptls.receive(conn)
  
    const res = check(eh_receive, {
      'received bytes': (r) => r.byteLength > 0,
    });
    errorRate.add(!res);
  }

  tcptls.closeConn(conn)
}
// 4. teardown code
export function teardown ()
{
}
