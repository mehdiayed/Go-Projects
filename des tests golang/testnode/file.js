const axios = require('axios');
const md5 = require('md5');

const readEgaugeData = async (egaugeJwT, deviceName) => {
  try {
    const auth = {
      headers: {
        Authorization: `Bearer ${egaugeJwT}`
      }
    };
    const config = {
      ...auth,
      params: {
        time: 'now',
        rate: '',
      },
    };
    const response = await axios.get(`https://${deviceName}.d.egauge.net/api/register`, config);
    return response.data;
  } catch (error) {
    console.log(error + " can't get data from egauge");
    return {};
  }
};

const egaugeLogin = async (deviceName, usr, pwd) => { 
  try {
    const unauthorizedResponse = await axios.get(`https://${deviceName}.d.egauge.net/api/auth/unauthorized`).catch(err => {
      return err.response?.data;
    });

    const rlm = unauthorizedResponse.rlm;
    const cnnc = Math.floor(Math.random() * 999999999).toString();

    const HA1 = md5(usr + ':' + rlm + ':' + pwd);
    const HA2 = md5(HA1 + ':' + unauthorizedResponse.nnc + ':' + cnnc);

    const loginBody = { rlm, usr: usr, cnnc, nnc: unauthorizedResponse.nnc, hash: HA2 };

    const response = await axios.post(`https://${deviceName}.d.egauge.net/api/auth/login`, loginBody);
    const jwt = response.data.jwt;
    return jwt;
  } catch (err) {
    console.log(err + "can't connect to egauge");
    return 'disconnected';
  }
};




/*

            here is the main program  

*/

const Device = "egauge67897"
const user = "owner"
const password = "000000"

const egaugeJwt = async () => await egaugeLogin(Device,user ,password );

(async () => {
  const jwt = await egaugeJwt();
  const data = await readEgaugeData(jwt, Device);
  console.log(data);
})();