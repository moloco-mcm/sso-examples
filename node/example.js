const { createHmac } = require("crypto");

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

// generateSignatureV100 creates the signature for v1.0.0
const generateSignatureV100 = (params, nonce, timestamp, secret) => {
  const {
    adAccountId,
    adAccountTitle,
    email,
    externalUserId,
    name,
    path,
    platformId,
    role,
    version,
  } = params;

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

  const concatenatedString = ssoParams.join("\n");
  return createHmac("sha256", secret)
    .update(concatenatedString)
    .digest("base64");
};

// generateSignatureV110 creates the signature for v1.1.0
// Parameters are in alphabetical order with new fields included
const generateSignatureV110 = (params, nonce, timestamp, secret) => {
  const {
    adAccountId,
    adAccountIds,
    adAccountTitle,
    adManagerAccountId,
    adManagerAccountTitle,
    email,
    externalUserId,
    name,
    path,
    platformId,
    role,
    version,
  } = params;

  // Join ad_account_ids array with semicolon separator if present
  const adAccountIdsStr =
    adAccountIds && adAccountIds.length > 0 ? adAccountIds.join(";") : "";

  // Parameters in alphabetical order
  const ssoParams = [
    adAccountId || "", // ad_account_id
    adAccountIdsStr, // ad_account_ids
    adAccountTitle || "", // ad_account_title
    adManagerAccountId || "", // ad_manager_account_id
    adManagerAccountTitle || "", // ad_manager_account_title
    email, // email
    externalUserId, // external_user_id
    name, // name
    nonce, // nonce
    path, // path
    platformId, // platform_id
    role || "", // role
    timestamp, // timestamp
    version, // version
  ];

  const concatenatedString = ssoParams.join("\n");
  return createHmac("sha256", secret)
    .update(concatenatedString)
    .digest("base64");
};

const createSignedRmpPortalUrl = (params) => {
  const {
    baseUrl,
    path,
    platformId,
    adAccountId,
    adAccountIds,
    adAccountTitle,
    adManagerAccountId,
    adManagerAccountTitle,
    name,
    email,
    role,
    externalUserId,
    secret,
    version,
    colorMode = "light",
    language,
  } = params;

  // unix timestamp in seconds
  const timestamp = `${Math.floor(Date.now() / 1000)}`;

  // create nonce
  const nonce = generateNonce(16);

  // Generate signature based on version
  let signature;
  if (version === "1.1.0") {
    signature = generateSignatureV110(params, nonce, timestamp, secret);
  } else {
    signature = generateSignatureV100(params, nonce, timestamp, secret);
  }

  // build url params
  const queryParams = {
    ad_account_id: adAccountId || "",
    ad_account_title: adAccountTitle || "",
    email,
    external_user_id: externalUserId,
    name,
    nonce,
    path,
    platform_id: platformId,
    role: role || "",
    timestamp,
    version,
    signature,
    "config:color_mode": colorMode,
    "config:language": language,
  };

  const queryString = new URLSearchParams(queryParams).toString();

  // Add v1.1.0 array fields separately (URLSearchParams doesn't handle arrays well)
  let finalQueryString = queryString;
  if (version === "1.1.0") {
    if (adAccountIds && adAccountIds.length > 0) {
      const adAccountIdsParams = adAccountIds
        .map((id) => `ad_account_ids=${encodeURIComponent(id)}`)
        .join("&");
      finalQueryString += `&${adAccountIdsParams}`;
    }
    if (adManagerAccountId) {
      finalQueryString += `&ad_manager_account_id=${encodeURIComponent(
        adManagerAccountId
      )}`;
    }
    if (adManagerAccountTitle) {
      finalQueryString += `&ad_manager_account_title=${encodeURIComponent(
        adManagerAccountTitle
      )}`;
    }
  }

  return `${baseUrl}/sso?${finalQueryString}`;
};

// exampleAdAccount demonstrates creating a signed URL for a single ad account.
// Works with AD_ACCOUNT_OWNER, AD_ACCOUNT_USER, or AD_ACCOUNT_VIEWER roles.
const exampleAdAccount = () => {
  const baseUrl = "https://{YOUR-RMP-PORTAL_URL}"; // Please use the url provided by your account manager
  const adAccountId = "my-ad-account-id";
  const adAccountTitle = "My Ad Account";
  const email = "test@example.com";
  const externalUserId = "user-id";
  const name = "Example User Name";
  const path = `/embed/sponsored-ads/cm/a/${adAccountId}`;
  const role = "AD_ACCOUNT_OWNER";
  const platformId = "RMP_PLATFORM_ID";
  const secret = "super-secret";
  const colorMode = "light"; // "light" | "dark" | "useDeviceSetting"
  const language = "en"; // "en" | "ko"
  const version = "1.1.0";

  return createSignedRmpPortalUrl({
    baseUrl,
    adAccountId,
    adAccountTitle,
    email,
    name,
    path,
    role,
    externalUserId,
    platformId,
    secret,
    version,
    colorMode,
    language,
  });
};

// exampleAdAccountAgency demonstrates creating a signed URL for multiple ad accounts with AD_ACCOUNT_AGENCY role (v1.1.0 feature).
// This allows an agency user to access multiple ad accounts.
const exampleAdAccountAgency = () => {
  const baseUrl = "https://{YOUR-RMP-PORTAL_URL}";
  const adAccountIds = ["account-1", "account-2", "account-3"];
  const email = "agency@example.com";
  const externalUserId = "agency-user-id";
  const name = "Agency User";
  const path = "/embed/sponsored-ads";
  const role = "AD_ACCOUNT_AGENCY";
  const platformId = "RMP_PLATFORM_ID";
  const secret = "super-secret";
  const colorMode = "light";
  const language = "en";
  const version = "1.1.0";

  return createSignedRmpPortalUrl({
    baseUrl,
    adAccountIds,
    email,
    name,
    path,
    role,
    externalUserId,
    platformId,
    secret,
    version,
    colorMode,
    language,
  });
};

// exampleAdManagerAccount demonstrates creating a signed URL for an ad manager account (v1.1.0 feature).
// Works with AD_MANAGER_ACCOUNT_OWNER or AD_MANAGER_ACCOUNT_USER roles.
const exampleAdManagerAccount = () => {
  const baseUrl = "https://{YOUR-RMP-PORTAL_URL}";
  const adManagerAccountId = "my-ad-manager-account-id";
  const adManagerAccountTitle = "My Ad Manager Account";
  const email = "manager@example.com";
  const externalUserId = "manager-user-id";
  const name = "Manager User";
  const path = `/embed/sponsored-ads/cm/ama/${adManagerAccountId}`;
  const role = "AD_MANAGER_ACCOUNT_OWNER";
  const platformId = "RMP_PLATFORM_ID";
  const secret = "super-secret";
  const colorMode = "light";
  const language = "en";
  const version = "1.1.0";

  return createSignedRmpPortalUrl({
    baseUrl,
    adManagerAccountId,
    adManagerAccountTitle,
    email,
    name,
    path,
    role,
    externalUserId,
    platformId,
    secret,
    version,
    colorMode,
    language,
  });
};

// Choose which example to run by uncommenting the desired line:
const signedUrl = exampleAdAccount();
// const signedUrl = exampleAdAccountAgency();
// const signedUrl = exampleAdManagerAccount();

console.log(signedUrl);
