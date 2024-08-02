import querystring from "querystring"
import axios, { AxiosResponse } from "axios"
import crypto from "crypto"
import { SPOTIFY_URL } from "../../env"

export class SpotifyAuthConfig {
  constructor(
    public authUrl: string,
    public tokenUrl: string,
    public redirectUri: string,
    public scope: string
  ) {}

  getAuthUrl(): string {
    return this.authUrl
  }

  getTokenUrl(): string {
    return this.tokenUrl
  }

  getRedirectUri(): string {
    return this.redirectUri
  }

  getScope(): string {
    return this.scope
  }
}

const generateRandomString = (length: number): string => {
  const possible =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
  const values = crypto.randomBytes(length)
  return Array.from(values)
    .map((byte) => possible[byte % possible.length])
    .join("")
}

export function redirectToSpotifyAuth(
  response: any,
  authConfig: SpotifyAuthConfig,
  clientId: string,
  randomStringLength: number
): void {
  const state = generateRandomString(randomStringLength)
  response.redirect(
    authConfig.getAuthUrl() +
      querystring.stringify({
        response_type: "code",
        client_id: clientId,
        scope: authConfig.getScope(),
        redirect_uri: authConfig.getRedirectUri(),
        state: state,
      })
  )
}

interface SpotifyAuthResponse {
  access_token: string
  refresh_token: string
}

const authUrl = `${SPOTIFY_URL}authorize`
const tokenUrl = `${SPOTIFY_URL}api/token`
const redirectUri =
  process.env.NEXT_PUBLIC_REDIRECT_URI || "http://localhost:3000/callback"
const scope =
  "playlist-modify-private playlist-modify-public playlist-read-private"

export const authConfig = new SpotifyAuthConfig(
  authUrl,
  tokenUrl,
  redirectUri,
  scope
)

export async function requestSpotifyToken(
  authCode: string,
  secretBase64: string
): Promise<SpotifyAuthResponse> {
  const tokenRequestOptions = {
    method: "post",
    url: authConfig.getTokenUrl(),
    data: querystring.stringify({
      code: authCode,
      redirect_uri: authConfig.getRedirectUri(),
      grant_type: "authorization_code",
    }),
    headers: {
      "content-type": "application/x-www-form-urlencoded",
      Authorization: "Basic " + secretBase64,
    },
  }

  return axios(tokenRequestOptions)
    .then((authResponse: AxiosResponse) => {
      return authResponse.data
    })
    .catch((error: any) => {
      throw new Error(`Could not get access token: ${error}`)
    })
}
