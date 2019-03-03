namespace brainz.Model
{
    using Newtonsoft.Json;

    public class CoverArtArchive
    {
        [JsonProperty("artwork")]
        public bool Artwork { get; set; }

        [JsonProperty("front")]
        public bool Front { get; set; }

        [JsonProperty("count")]
        public long Count { get; set; }

        [JsonProperty("back")]
        public bool Back { get; set; }

        [JsonProperty("darkened")]
        public bool Darkened { get; set; }
    }
}
