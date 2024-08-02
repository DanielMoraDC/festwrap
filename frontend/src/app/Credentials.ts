export class SpotifyCredentials {
  private clientId: string
  private secret: string

  constructor(clientId: string, secret: string) {
    this.clientId = clientId
    this.secret = secret
  }

  getClientId() {
    return this.clientId
  }

  getBase64Secret() {
    return Buffer.from(this.clientId + ":" + this.secret).toString("base64")
  }
}
