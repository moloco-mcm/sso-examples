import java.net.URLEncoder;
import java.security.SecureRandom;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Base64;
import java.util.List;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class EmbeddedCampaignManagerExample {
  private static final String HMAC_SHA256_ALGORITHM = "HmacSHA256";
  private static final String UTF_8_CHARSET = "UTF-8";
  // Please use the url provided by your account manager
  private static final String RMP_PORTAL_BASE_URL = "https://{YOUR-RMP-PORTAL_URL}";

  public static void main(String[] args) {
    // Choose which example to run by uncommenting the desired line:
    try {
      String url = exampleAdAccount();
      // String url = exampleAdAccountAgency();
      // String url = exampleAdManagerAccount();

      System.out.println(url);
    } catch (Exception e) {
      System.out.println(e);
    }
  }

  // exampleAdAccount demonstrates creating a signed URL for a single ad account.
  // Works with AD_ACCOUNT_OWNER, AD_ACCOUNT_USER, or AD_ACCOUNT_VIEWER roles.
  public static String exampleAdAccount() throws Exception {
    String baseUrl = RMP_PORTAL_BASE_URL;
    String platformId = "RMP_PLATFORM_ID";
    String adAccountId = "my-ad-account-id";
    String adAccountTitle = "My Ad Account";
    String path = "/embed/sponsored-ads/cm/a/" + adAccountId;
    String name = "Example User Name";
    String email = "test@example.com";
    String role = "AD_ACCOUNT_OWNER";
    String externalUserId = "user-id";
    String secret = "super-secret";
    String version = "1.1.0";
    String colorMode = "light"; // "light" | "dark" | "useDeviceSetting"
    String language = "en"; // "en" | "ko"

    return createSignedRmpPortalUrl(baseUrl, path, platformId, adAccountId, null, adAccountTitle, null, null,
        name, email, role, externalUserId, secret, version, colorMode, language);
  }

  // exampleAdAccountAgency demonstrates creating a signed URL for multiple ad
  // accounts with AD_ACCOUNT_AGENCY role.
  // This allows an agency user to access multiple ad accounts.
  public static String exampleAdAccountAgency() throws Exception {
    String baseUrl = RMP_PORTAL_BASE_URL;
    String platformId = "RMP_PLATFORM_ID";
    String[] adAccountIds = new String[] { "account-1", "account-2", "account-3" };
    String path = "/embed/sponsored-ads";
    String name = "Agency User";
    String email = "agency@example.com";
    String role = "AD_ACCOUNT_AGENCY";
    String externalUserId = "agency-user-id";
    String secret = "super-secret";
    String version = "1.1.0";
    String colorMode = "light";
    String language = "en";

    return createSignedRmpPortalUrl(baseUrl, path, platformId, "", adAccountIds, "", null, null,
        name, email, role, externalUserId, secret, version, colorMode, language);
  }

  // exampleAdManagerAccount demonstrates creating a signed URL for an ad manager
  // account.
  // Works with AD_MANAGER_ACCOUNT_OWNER or AD_MANAGER_ACCOUNT_USER roles.
  public static String exampleAdManagerAccount() throws Exception {
    String baseUrl = RMP_PORTAL_BASE_URL;
    String platformId = "RMP_PLATFORM_ID";
    String adManagerAccountId = "my-ad-manager-account-id";
    String adManagerAccountTitle = "My Ad Manager Account";
    String path = "/embed/sponsored-ads/cm/ama/" + adManagerAccountId;
    String name = "Manager User";
    String email = "manager@example.com";
    String role = "AD_MANAGER_ACCOUNT_OWNER";
    String externalUserId = "manager-user-id";
    String secret = "super-secret";
    String version = "1.1.0";
    String colorMode = "light";
    String language = "en";

    return createSignedRmpPortalUrl(baseUrl, path, platformId, "", null, "", adManagerAccountId,
        adManagerAccountTitle, name, email, role, externalUserId, secret, version, colorMode, language);
  }

  public static String buildQueryParam(String name, String value) throws Exception {
    return name + "=" + URLEncoder.encode(value, UTF_8_CHARSET);
  };

  public static String createSignedRmpPortalUrl(String baseUrl, String path, String platformId, String adAccountId,
      String[] adAccountIds, String adAccountTitle, String adManagerAccountId, String adManagerAccountTitle,
      String name, String email, String role, String externalUserId, String secret, String version, String colorMode,
      String language) throws Exception {
    // unix timestamp in seconds
    String timestamp = Long.toString(System.currentTimeMillis() / 1000L);

    // create nonce
    SecureRandom random = new SecureRandom();
    byte bytes[] = new byte[20];
    random.nextBytes(bytes);
    String nonce = Base64.getEncoder().encodeToString(bytes);

    // Generate signature
    String signature = generateSignature(adAccountId, adAccountIds, adAccountTitle, adManagerAccountId,
        adManagerAccountTitle, email, externalUserId, name, nonce, path, platformId, role, timestamp, version, secret);

    // construct final url with query params
    List<String> queryParamsList = new ArrayList<>();
    queryParamsList.add(buildQueryParam("ad_account_id", adAccountId != null ? adAccountId : ""));
    queryParamsList.add(buildQueryParam("ad_account_title", adAccountTitle != null ? adAccountTitle : ""));
    queryParamsList.add(buildQueryParam("email", email));
    queryParamsList.add(buildQueryParam("external_user_id", externalUserId));
    queryParamsList.add(buildQueryParam("name", name));
    queryParamsList.add(buildQueryParam("nonce", nonce));
    queryParamsList.add(buildQueryParam("path", path));
    queryParamsList.add(buildQueryParam("platform_id", platformId));
    queryParamsList.add(buildQueryParam("role", role));
    queryParamsList.add(buildQueryParam("timestamp", timestamp));
    queryParamsList.add(buildQueryParam("version", version));
    queryParamsList.add(buildQueryParam("signature", signature));

    // Add optional fields if present
    if (adAccountIds != null && adAccountIds.length > 0) {
      for (String id : adAccountIds) {
        queryParamsList.add(buildQueryParam("ad_account_ids", id));
      }
    }
    if (adManagerAccountId != null && !adManagerAccountId.isEmpty()) {
      queryParamsList.add(buildQueryParam("ad_manager_account_id", adManagerAccountId));
    }
    if (adManagerAccountTitle != null && !adManagerAccountTitle.isEmpty()) {
      queryParamsList.add(buildQueryParam("ad_manager_account_title", adManagerAccountTitle));
    }

    queryParamsList.add(buildQueryParam("config:color_mode", colorMode));
    queryParamsList.add(buildQueryParam("config:language", language));

    String signedUrl = baseUrl + "/sso?" + String.join("&", queryParamsList);

    return signedUrl;
  }

  // generateSignature creates the HMAC-SHA256 signature for the SSO URL.
  // Parameters are in alphabetical order.
  private static String generateSignature(String adAccountId, String[] adAccountIds, String adAccountTitle,
      String adManagerAccountId, String adManagerAccountTitle, String email, String externalUserId, String name,
      String nonce, String path, String platformId, String role, String timestamp, String version, String secret)
      throws Exception {
    // Join ad_account_ids array with semicolon separator if present
    String adAccountIdsStr = "";
    if (adAccountIds != null && adAccountIds.length > 0) {
      adAccountIdsStr = String.join(";", adAccountIds);
    }

    // Parameters in alphabetical order
    String[] paramArray = new String[] {
        adAccountId != null ? adAccountId : "", // ad_account_id
        adAccountIdsStr, // ad_account_ids
        adAccountTitle != null ? adAccountTitle : "", // ad_account_title
        adManagerAccountId != null ? adManagerAccountId : "", // ad_manager_account_id
        adManagerAccountTitle != null ? adManagerAccountTitle : "", // ad_manager_account_title
        email, // email
        externalUserId, // external_user_id
        name, // name
        nonce, // nonce
        path, // path
        platformId, // platform_id
        role, // role
        timestamp, // timestamp
        version // version
    };
    String paramString = String.join("\n", Arrays.asList(paramArray));

    byte[] keyBytes = secret.getBytes();
    SecretKeySpec signingKey = new SecretKeySpec(keyBytes, HMAC_SHA256_ALGORITHM);
    Mac mac = Mac.getInstance(HMAC_SHA256_ALGORITHM);
    mac.init(signingKey);
    byte[] rawHmac = Base64.getEncoder().encode(mac.doFinal(paramString.getBytes(UTF_8_CHARSET)));
    return new String(rawHmac, UTF_8_CHARSET);
  }
}
