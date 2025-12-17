const { createHmac } = require("crypto");

const createSignedRmpPortalUrl = (params) => {
  const {
    baseUrl,
    path,
    platformId,
    adAccountId,
    adAccountTitle,
    name,
    email,
    role,
    externalUserId,
    nonce,
    secret,
    version,
    colorMode = "light",
    language,
  } = params;

  // unix timestamp in seconds
  const timestamp = `${Math.floor(Date.now() / 1000)}`;

  // gather params in an alphabetical order
  const ssoParams = [
    adAccountId || "",
    adAccountTitle || "",
    email,
    externalUserId,
    name,
    nonce,
    path,
    platformId,
    role || "",
    timestamp,
    version,
  ];

  // concatenate the values with line break “\n”
  const concatenatedString = ssoParams.join("\n");

  // create a signature
  const signature = createHmac("sha256", secret)
    .update(concatenatedString)
    .digest("base64");

  // build url params
  const queryString = new URLSearchParams({
    ad_account_id: adAccountId,
    ad_account_title: adAccountTitle,
    email,
    external_user_id: externalUserId,
    name,
    nonce,
    path,
    platform_id: platformId,
    role,
    timestamp,
    version,
    signature,
    "config:color_mode": colorMode,
    "config:language": language,
  }).toString();

  return `${baseUrl}/sso?${queryString}`;
};

const generateNonce = (length) => {
  const possibleCharacters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  return Array(length)
    .fill(0)
    .map(() =>
      possibleCharacters.charAt(
        Math.floor(Math.random() * possibleCharacters.length)
      )
    )
    .join("");
};

const baseUrl = "https://{YOUR-RMP-PORTAL_URL}"; // Please use the url provided by your account manager
const adAccountId = "ad-account-id";
const adAccountTitle = "My Ad Account";
const email = "test@example.com";
const externalUserId = "user-id";
const name = "Example User Name";
const nonce = generateNonce(16);
const path = `/embed/sponsored-ads/cm/a/${adAccountId}`;
const role = "AD_ACCOUNT_OWNER";
const platformId = "RMP_PLATFORM_ID";
const secret = "super-secret";
const colorMode = "light"; // "light" | "dark" | "useDeviceSetting"
const language = "en"; // "en" | "ko"
const version = "1.0.0";

const signedUrl = createSignedRmpPortalUrl({
  baseUrl,
  adAccountId,
  adAccountTitle,
  email,
  name,
  nonce,
  path,
  role,
  externalUserId,
  platformId,
  secret,
  version,
  colorMode,
  language,
});

console.log(signedUrl);
