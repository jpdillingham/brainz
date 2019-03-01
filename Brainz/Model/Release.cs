namespace brainz.Model
{
    using System;
    using Newtonsoft.Json;

    public partial class Release
    {
        [JsonProperty("packaging-id")]
        public Guid? PackagingId { get; set; }

        [JsonProperty("asin")]
        public string Asin { get; set; }

        [JsonProperty("status-id")]
        public Guid? StatusId { get; set; }

        [JsonProperty("disambiguation")]
        public string Disambiguation { get; set; }

        [JsonProperty("date")]
        public string Date { get; set; }

        [JsonProperty("packaging")]
        public string Packaging { get; set; }

        [JsonProperty("status")]
        public string Status { get; set; }

        [JsonProperty("release-events")]
        public ReleaseEvent[] ReleaseEvents { get; set; }

        [JsonProperty("cover-art-archive")]
        public CoverArtArchive CoverArtArchive { get; set; }

        [JsonProperty("text-representation")]
        public TextRepresentation TextRepresentation { get; set; }

        [JsonProperty("quality")]
        public string Quality { get; set; }

        [JsonProperty("title")]
        public string Title { get; set; }

        [JsonProperty("country")]
        public string Country { get; set; }

        [JsonProperty("id")]
        public Guid Id { get; set; }

        [JsonProperty("media")]
        public Media[] Media { get; set; }

        [JsonProperty("barcode")]
        public string Barcode { get; set; }
    }
}