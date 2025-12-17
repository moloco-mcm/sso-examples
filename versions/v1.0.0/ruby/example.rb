require 'securerandom'
require 'base64'
require 'openssl'
require 'uri'

def create_signed_rmp_portal_url(args)
  base_url = args[:base_url]
  ad_account_id = args[:ad_account_id]
  ad_account_title = args[:ad_account_title]
  email = args[:email]
  external_user_id = args[:external_user_id]
  name = args[:name]
  path = args[:path]
  platform_id = args[:platform_id]
  role = args[:role]
  version = args[:version]
  color_mode = args[:color_mode]
  language = args[:language]
  secret = args[:secret]
  timestamp = Time.now.to_i.to_s
  nonce = SecureRandom.hex(10).to_s

  param_array = [ad_account_id, ad_account_title, email, external_user_id, name, nonce, path,
                 platform_id, role, timestamp, version]
  param_string = param_array.join("\n")

  signature = Base64.encode64(
    OpenSSL::HMAC.digest(
      OpenSSL::Digest.new('sha256'),
      secret,
      param_string.force_encoding('utf-8')
    )
  ).strip

  query_params = {
    ad_account_id: ad_account_id,
    ad_account_title: ad_account_title,
    email: email,
    external_user_id: external_user_id,
    name: name,
    nonce: nonce,
    path: path,
    platform_id: platform_id,
    role: role,
    timestamp: timestamp,
    version: version,
    signature: signature,
    'config:color_mode': color_mode,
    'config:language': language
  }

  query_string = URI.encode_www_form(query_params)

  base_url + '/sso?' + query_string
end

ad_account_id = 'my_ad_account_id'

args = {
  base_url: 'https://{YOUR-RMP-PORTAL_URL}', # Please use the url provided by your account manager
  ad_account_id: ad_account_id,
  path: '/embed/sponsored-ads/cm/a/' + ad_account_id,
  platform_id: 'RMP_PLATFORM_ID',
  ad_account_title: 'My Ad Account',
  name: 'Example User Name',
  email: 'test@example.com',
  role: 'AD_ACCOUNT_OWNER',
  external_user_id: 'user-id',
  secret: 'super-secret',
  version: '1.0.0',
  color_mode: 'light',
  language: 'en'
}

signed_url = create_signed_rmp_portal_url(args)

puts(signed_url)
