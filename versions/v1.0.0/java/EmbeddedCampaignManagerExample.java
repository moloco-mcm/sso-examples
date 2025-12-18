import java.net.URLEncoder;
import java.security.SecureRandom;
import java.util.Arrays;
import java.util.Base64;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class EmbeddedCampaignManagerExample {
  private static final String HMAC_SHA256_ALGORITHM = "HmacSHA256";
  private static final String UTF_8_CHARSET = "UTF-8";
  // Please use the url provided by your account manager
  private static final String RMP_PORTAL_BASE_URL = "https://{YOUR-RMP-PORTAL_URL}";

  public static void main(String[] args) {
    String baseUrl = RMP_PORTAL_BASE_URL;
    String platformId = "RMP_PLATFORM_ID";
    String adAccountId = "ad-account-id";
    String adAccountTitle = "My Ad Account";
    String path = "/embed/sponsored-ads/cm/a/" + adAccountId;
    String name = "Example User Name";
    String email = "test@example.com";
    String role = "AD_ACCOUNT_OWNER";
    String externalUserId = "user-id";
    String secret = "super-secret";
    String version = "1.0.0";
    String colorMode = "light"; // "light" | "dark" | "useDeviceSetting"
    String language = "en"; // "en" | "ko"

    try {
      String url = createSignedRmpPortalUrl(baseUrl, path, platformId, adAccountId, adAccountTitle, name, email, role,
          externalUserId, secret, version, colorMode, language);
      System.out.println(url);

    } catch (Exception e) {
      System.out.println(e);
    }
  }

  public static String buildQueryParam(String name, String value) throws Exception {
    return name + "=" + URLEncoder.encode(value, UTF_8_CHARSET);
  };

  public static String createSignedRmpPortalUrl(String baseUrl, String path, String platformId, String adAccountId,
      String adAccountTitle, String name, String email, String role, String externalUserId, String secret,
      String version, String colorMode, String language) throws Exception {
    // unix timestamp in seconds
    String timestamp = Long.toString(System.currentTimeMillis() / 1000L);

    // create nonce
    SecureRandom random = new SecureRandom();
    byte bytes[] = new byte[20];
    random.nextBytes(bytes);
    String nonce = Base64.getEncoder().encodeToString(bytes);

    // gather params in an alphabetical order
    String[] paramArray = new String[] { adAccountId, adAccountTitle, email, externalUserId, name, nonce, path,
        platformId, role, timestamp, version };
    String paramString = String.join("\n", Arrays.asList(paramArray));

    // create signature
    byte[] keyBytes = secret.getBytes();
    SecretKeySpec signingKey = new SecretKeySpec(keyBytes, HMAC_SHA256_ALGORITHM);
    Mac mac = Mac.getInstance(HMAC_SHA256_ALGORITHM);
    mac.init(signingKey);
    byte[] rawHmac = Base64.getEncoder().encode(mac.doFinal(paramString.getBytes(UTF_8_CHARSET)));
    String signature = new String(rawHmac, UTF_8_CHARSET);

    // construct final url
    String[] queryParams = new String[] { buildQueryParam("ad_account_id", adAccountId),
        buildQueryParam("ad_account_title", adAccountTitle), buildQueryParam("email", email),
        buildQueryParam("external_user_id", externalUserId), buildQueryParam("name", name),
        buildQueryParam("nonce", nonce), buildQueryParam("path", path), buildQueryParam("platform_id", platformId),
        buildQueryParam("role", role), buildQueryParam("timestamp", timestamp), buildQueryParam("version", version),
        buildQueryParam("signature", signature), buildQueryParam("config:color_mode", colorMode),
        buildQueryParam("config:language", language) };

    String signedUrl = baseUrl + "/sso?" + String.join("&", Arrays.asList(queryParams));

    return signedUrl;
  }
}
