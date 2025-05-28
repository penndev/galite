// Author: Penndev
// Remark: 
//  对请求网址进行签名，防止篡改并添加资源访问超时控制。
//  对请求数据进行加密，并解析返回的加密数据。


/**
 * base64 URL编码
 * https://datatracker.ietf.org/doc/html/rfc4648#section-5
 * @param {Uint8Array} str 
 * @returns base64url string
 */
function base64URLEncode(bytes) {
  const base64urlChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_';
  let output = '';
  let i;

  for (i = 0; i + 2 < bytes.length; i += 3) {
    const triplet = (bytes[i] << 16) | (bytes[i + 1] << 8) | bytes[i + 2];
    output += base64urlChars[(triplet >> 18) & 0x3F];
    output += base64urlChars[(triplet >> 12) & 0x3F];
    output += base64urlChars[(triplet >> 6) & 0x3F];
    output += base64urlChars[triplet & 0x3F];
  }

  const remaining = bytes.length - i;
  if (remaining === 1) {
    const triplet = bytes[i] << 16;
    output += base64urlChars[(triplet >> 18) & 0x3F];
    output += base64urlChars[(triplet >> 12) & 0x3F];
    // No padding in base64url
  } else if (remaining === 2) {
    const triplet = (bytes[i] << 16) | (bytes[i + 1] << 8);
    output += base64urlChars[(triplet >> 18) & 0x3F];
    output += base64urlChars[(triplet >> 12) & 0x3F];
    output += base64urlChars[(triplet >> 6) & 0x3F];
  }

  return output;
}

/**
 * base64 URL解码
 * @param {string} str base64url字符串
 * @returns {Uint8Array}
 */
function base64URLDecode(str) {
  // 替换URL安全字符
  str = str.replace(/-/g, '+').replace(/_/g, '/');
  // 补齐padding
  while (str.length % 4) {
    str += '=';
  }
  // 手动解码base64字符串为Uint8Array
  const base64Chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
  let buffer = [];
  let bits = 0, bitLength = 0;
  for (let i = 0; i < str.length; i++) {
    const c = str[i];
    if (c === '=') break;
    const val = base64Chars.indexOf(c);
    if (val === -1) continue;
    bits = (bits << 6) | val;
    bitLength += 6;
    if (bitLength >= 8) {
      bitLength -= 8;
      buffer.push((bits >> bitLength) & 0xFF);
    }
  }
  return new Uint8Array(buffer);
}

/**
 * 使用HMAC对消息进行签名，并返回base64 URL编码的签名。
 * https://developer.mozilla.org/zh-CN/docs/Web/API/SubtleCrypto/importKey
 * @param {string} msg 签名内容
 * @param {string} key 密钥
 * @param {string} algorithm 算法类型
 * @returns Uint8Array
 */
const signHMAC = async (msg, key, algorithm = "SHA-256") => {
  const encoder = new TextEncoder();
  const msgData = encoder.encode(msg);
  const keyRaw = encoder.encode(key);
  const cryptoKey = await crypto.subtle.importKey(
    "raw", keyRaw,
    { name: "HMAC", hash: { name: algorithm } },
    false, ["sign"]
  );
  const signature = await crypto.subtle.sign("HMAC", cryptoKey, msgData);
  return new Uint8Array(signature);
}

/**
 * 使用 AES-GCM 加密数据（浏览器环境）
 * @param {Uint8Array} msg - 明文
 * @param {string} key - 原始密钥（16/24/32字节）
 * @param {Uint8Array} iv - 初始化向量（12字节推荐）
 * @returns {Promise<Uint8Array>} - 加密后的密文
 */
const encryptAesGCM = async (msg, key, iv) => {
  const encoder = new TextEncoder();
  const msgRaw = encoder.encode(msg);
  const keyByte = encoder.encode(key);
  const keyRaw = await crypto.subtle.digest("SHA-256", keyByte)
  const cryptoKey = await crypto.subtle.importKey(
    "raw", keyRaw,
    { name: "AES-GCM" },
    false, ["encrypt"]
  );
  const encrypted = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv: iv },
    cryptoKey, msgRaw
  );
  return new Uint8Array(encrypted);
}

/**
 * 使用 AES-GCM 解密数据（浏览器环境）
 * @param {Uint8Array} msg - 密文
 * @param {string} key - 原始密钥（16/24/32字节）
 * @param {Uint8Array} iv - 初始化向量（12字节推荐）
 * @returns {Promise<Uint8Array>} - 解密后的明文
 */
const decryptAesGCM = async (msg, key, iv) => {
  const encoder = new TextEncoder();
  const keyByte = encoder.encode(key);
  const keyRaw = await crypto.subtle.digest("SHA-256", keyByte);

  const cryptoKey = await crypto.subtle.importKey(
    "raw", keyRaw,
    { name: "AES-GCM" },
    false, ["decrypt"]
  );
  // 解密
  const decrypted = await crypto.subtle.decrypt(
    { name: "AES-GCM", iv: iv, },
    cryptoKey, msg
  );
  return new Uint8Array(decrypted);
}



/**
 * 签名部分
 * - 验证并设置过期时间 expired
 * - 添加签名信息 sign
 */
const sign = async (url, key, expire) => {
  const originURL = new URL(url);
  if (originURL.hash) {
    throw new Error(`URL contains a hash fragment: ${originURL.toString()}`);
  }
  if (originURL.searchParams.has("expired")) {
    throw new Error(`URL expired: ${originURL.toString()}`);
  }
  const expired = Math.floor(new Date().getTime() / 1000) + expire; // 30秒后过期
  originURL.searchParams.set("expired", expired);
  if (originURL.searchParams.has("sign")) {
    throw new Error(`url signature is already set: ${originURL.toString()}`);
  }
  const signMsg = originURL.pathname + originURL.search
  const signUint8Array = await signHMAC(signMsg, key, "SHA-256");
  const sign = base64URLEncode(signUint8Array);
  const signUrl = originURL.toString() + `&sign=${sign}`;
  // 输出: https://example.com/api?param=value&expired=100&sign=LQUFO-b0lrXwMHmd9ERoLBxQld7iMtJMX4VD3V8e0Qs
  return signUrl;
}


/**
 * 请求加密部分
 * - 请求加密
 * - 返回解密
 */
const encryptFetch = async (signUrl, method, body, key) => {
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const response = await fetch(signUrl, {
    method: method,
    headers: {
      "Content-Type": "application/x-buffer",
      "X-Iv": base64URLEncode(iv),
    },
    body: await encryptAesGCM(body, key, iv),
  })

  var respBody;
  if (response.headers.get("Content-Type") == "application/x-buffer") {
    const iv = base64URLDecode(response.headers.get("X-Iv"));
    if (iv.length == 12) {
      const respMessage = new Uint8Array(await response.arrayBuffer());
      const respBodyBuffer = await decryptAesGCM(respMessage, key, iv);
      respBody = new TextDecoder().decode(respBodyBuffer)
    }
  } else {
    respBody = await response.text()
  }
  return respBody
}

const key = "secret"  //加密密钥
const body = `{"message":"hello example.js"}`;
const url = "http://127.0.0.1:8000/encrypt?param=value"

const signUrl = await sign(url, key, 30);
const respBody = await encryptFetch(signUrl, "POST", body, key)
if (respBody != body) {
  throw new Error(`Response body does not match request body. Expected: ${body}, Received: ${respBody}`);
}




// 基准测试
// 普通发送 10000 次 总耗时:4150.405ms  平均耗时：0.4150405ms
// 普通发送 50000 次 总耗时:18259.5003ms  平均耗时：0.365190006ms

// 加密发送 10000 次 总耗时:8042.4536ms  平均耗时：0.80424536ms
// 加密发送 50000 次 总耗时:34866.3655ms  平均耗时：0.69732731ms




