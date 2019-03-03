namespace brainz.Responses
{
    using Newtonsoft.Json;
    using brainz.Model;

    public partial class ReleaseResponse
    {
        [JsonProperty("release-offset")]
        public long ReleaseOffset { get; set; }

        [JsonProperty("releases")]
        public Release[] Releases { get; set; }

        [JsonProperty("release-count")]
        public long ReleaseCount { get; set; }
    }
}
