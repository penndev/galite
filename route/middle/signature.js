// Author: Penndev
// Remark: 对请求网址进行签名，防止篡改并添加资源访问超时控制。


/**
 * base64 URL编码
 * https://datatracker.ietf.org/doc/html/rfc4648#section-5
 * @param {string} str 
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
    const keyData = encoder.encode(key);
    const cryptoKey = await crypto.subtle.importKey(
        "raw", keyData,
        { name: "HMAC", hash: { name: algorithm } },
        false, ["sign"]
    );
    const signature = await crypto.subtle.sign("HMAC", cryptoKey, msgData);
    return new Uint8Array(signature);
}

const originURL = new URL("http://127.0.0.1:8000/ping?param=value");
if (originURL.hash) {
    throw new Error(`URL contains a hash fragment: ${originURL.toString()}`);;
}

// 验证并设置过期时间
if (originURL.searchParams.has("expired")) {
    throw new Error(`URL expired: ${originURL.toString()}`);
}
const expired = Math.floor(new Date().getTime() / 1000) + 30; // 30秒后过期
originURL.searchParams.set("expired", expired);

// 添加签名
if (originURL.searchParams.has("sign")) {
    throw new Error(`url signature is already set: ${originURL.toString()}`);
}
const signMsg = originURL.pathname + originURL.search
const signUint8Array = await signHMAC(signMsg, "secret", "SHA-256");
const sign = base64URLEncode(signUint8Array);
console.log(originURL.toString() + `&sign=${sign}`);
// 输出: https://example.com/api?param=value&expired=100&sign=LQUFO-b0lrXwMHmd9ERoLBxQld7iMtJMX4VD3V8e0Qs