"use client"

import { useRouter, useSearchParams } from "next/navigation"
import { useEffect } from "react"
import { SPOTIFY_CLIENT_ID, SPOTIFY_SECRET } from "../../../env"
import { requestSpotifyToken } from "../Authorization"
import { SpotifyCredentials } from "../Credentials"

const Callback = () => {
  const searchParams = useSearchParams()
  const router = useRouter()

  const authCode = searchParams.get("code")
  const state = searchParams.get("state")

  useEffect(() => {
    async function requestAndStoreToken() {
      if (authCode && SPOTIFY_CLIENT_ID && SPOTIFY_SECRET) {
        const spotifyCredentials = new SpotifyCredentials(
          SPOTIFY_CLIENT_ID,
          SPOTIFY_SECRET
        )
        const response = await requestSpotifyToken(
          authCode,
          spotifyCredentials.getBase64Secret()
        )

        localStorage.setItem("access_token", response.access_token)
        localStorage.setItem("refresh_token", response.refresh_token)
      }
    }

    requestAndStoreToken()
    router.push("/")
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return <div>Redirecting...</div>
}

export default Callback
