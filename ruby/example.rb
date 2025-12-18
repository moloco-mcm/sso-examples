require 'securerandom'
require 'base64'
require 'openssl'
require 'uri'

# generateSignature creates the HMAC-SHA256 signature for the SSO URL.
# Parameters are in alphabetical order.
def generate_signature(args, nonce, timestamp, secret, version)
  # Join ad_account_ids array with semicolon separator if present
  ad_account_ids_str = ''
  if args[:ad_account_ids] && args[:ad_account_ids].length > 0
    ad_account_ids_str = args[:ad_account_ids].join(';')
  end

  # Parameters in alphabetical order
  param_array = [
    args[:ad_account_id] || '',                    # ad_account_id
    ad_account_ids_str,                            # ad_account_ids
    args[:ad_account_title] || '',                 # ad_account_title
    args[:ad_manager_account_id] || '',            # ad_manager_account_id
    args[:ad_manager_account_title] || '',         # ad_manager_account_title
    args[:email],                                  # email
    args[:external_user_id],                       # external_user_id
    args[:name],                                   # name
    nonce,                                         # nonce
    args[:path],                                   # path
    args[:platform_id],                            # platform_id
    args[:role] || '',                             # role
    timestamp,                                     # timestamp
    version                                        # version
  ]
  param_string = param_array.join("\n")

  Base64.encode64(
    OpenSSL::HMAC.digest(
      OpenSSL::Digest.new('sha256'),
      secret,
      param_string.force_encoding('utf-8')
    )
  ).strip
end

def create_signed_rmp_portal_url(args)
  base_url = args[:base_url]
  secret = args[:secret]
  version = args[:version]
  timestamp = Time.now.to_i.to_s
  nonce = SecureRandom.hex(10).to_s

  # Generate signature
  signature = generate_signature(args, nonce, timestamp, secret, version)

  query_params = {
    ad_account_id: args[:ad_account_id] || '',
    ad_account_title: args[:ad_account_title] || '',
    email: args[:email],
    external_user_id: args[:external_user_id],
    name: args[:name],
    nonce: nonce,
    path: args[:path],
    platform_id: args[:platform_id],
    role: args[:role] || '',
    timestamp: timestamp,
    version: version,
    signature: signature,
    'config:color_mode': args[:color_mode],
    'config:language': args[:language]
  }

  query_string = URI.encode_www_form(query_params)

  # Add optional fields if present
  if args[:ad_account_ids] && args[:ad_account_ids].length > 0
    ad_account_ids_params = args[:ad_account_ids].map { |id| "ad_account_ids=#{URI.encode_www_form_component(id)}" }.join('&')
    query_string += "&#{ad_account_ids_params}"
  end
  if args[:ad_manager_account_id] && !args[:ad_manager_account_id].empty?
    query_string += "&ad_manager_account_id=#{URI.encode_www_form_component(args[:ad_manager_account_id])}"
  end
  if args[:ad_manager_account_title] && !args[:ad_manager_account_title].empty?
    query_string += "&ad_manager_account_title=#{URI.encode_www_form_component(args[:ad_manager_account_title])}"
  end

  base_url + '/sso?' + query_string
end

# exampleAdAccount demonstrates creating a signed URL for a single ad account.
# Works with AD_ACCOUNT_OWNER, AD_ACCOUNT_USER, or AD_ACCOUNT_VIEWER roles.
def example_ad_account
  ad_account_id = 'my-ad-account-id'

  args = {
    base_url: 'https://{YOUR-RMP-PORTAL_URL}', # Please use the url provided by your account manager
    ad_account_id: ad_account_id,
    ad_account_title: 'My Ad Account',
    email: 'test@example.com',
    external_user_id: 'user-id',
    name: 'Example User Name',
    path: '/embed/sponsored-ads/cm/a/' + ad_account_id,
    platform_id: 'RMP_PLATFORM_ID',
    role: 'AD_ACCOUNT_OWNER',
    secret: 'super-secret',
    version: '1.1.0',
    color_mode: 'light',
    language: 'en'
  }

  create_signed_rmp_portal_url(args)
end

# exampleAdAccountAgency demonstrates creating a signed URL for multiple ad accounts with AD_ACCOUNT_AGENCY role.
# This allows an agency user to access multiple ad accounts.
def example_ad_account_agency
  args = {
    base_url: 'https://{YOUR-RMP-PORTAL_URL}',
    ad_account_ids: ['account-1', 'account-2', 'account-3'],
    email: 'agency@example.com',
    external_user_id: 'agency-user-id',
    name: 'Agency User',
    path: '/embed/sponsored-ads',
    platform_id: 'RMP_PLATFORM_ID',
    role: 'AD_ACCOUNT_AGENCY',
    secret: 'super-secret',
    version: '1.1.0',
    color_mode: 'light',
    language: 'en'
  }

  create_signed_rmp_portal_url(args)
end

# exampleAdManagerAccount demonstrates creating a signed URL for an ad manager account.
# Works with AD_MANAGER_ACCOUNT_OWNER or AD_MANAGER_ACCOUNT_USER roles.
def example_ad_manager_account
  ad_manager_account_id = 'my-ad-manager-account-id'

  args = {
    base_url: 'https://{YOUR-RMP-PORTAL_URL}',
    ad_manager_account_id: ad_manager_account_id,
    ad_manager_account_title: 'My Ad Manager Account',
    email: 'manager@example.com',
    external_user_id: 'manager-user-id',
    name: 'Manager User',
    path: '/embed/sponsored-ads/cm/ama/' + ad_manager_account_id,
    platform_id: 'RMP_PLATFORM_ID',
    role: 'AD_MANAGER_ACCOUNT_OWNER',
    secret: 'super-secret',
    version: '1.1.0',
    color_mode: 'light',
    language: 'en'
  }

  create_signed_rmp_portal_url(args)
end

# Choose which example to run by uncommenting the desired line:
signed_url = example_ad_account
# signed_url = example_ad_account_agency
# signed_url = example_ad_manager_account

puts(signed_url)
