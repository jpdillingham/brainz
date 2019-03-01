namespace Brainz.Responses
{
    using System;
    using Brainz.Model;
    using Newtonsoft.Json;

    public class RecordingResponse
    {
        [JsonProperty("media")]
        public Media[] Media { get; set; }

        [JsonProperty("packaging-id")]
        public Guid PackagingId { get; set; }

        [JsonProperty("quality")]
        public string Quality { get; set; }

        [JsonProperty("country")]
        public string Country { get; set; }

        [JsonProperty("title")]
        public string Title { get; set; }

        [JsonProperty("disambiguation")]
        public string Disambiguation { get; set; }

        [JsonProperty("packaging")]
        public string Packaging { get; set; }

        [JsonProperty("release-events")]
        public ReleaseEvent[] ReleaseEvents { get; set; }

        [JsonProperty("date")]
        public string Date { get; set; }

        [JsonProperty("cover-art-archive")]
        public CoverArtArchive CoverArtArchive { get; set; }

        [JsonProperty("id")]
        public Guid Id { get; set; }

        [JsonProperty("asin")]
        public object Asin { get; set; }

        [JsonProperty("text-representation")]
        public TextRepresentation TextRepresentation { get; set; }

        [JsonProperty("status")]
        public string Status { get; set; }

        [JsonProperty("barcode")]
        public string Barcode { get; set; }

        [JsonProperty("status-id")]
        public Guid StatusId { get; set; }
    }
}
