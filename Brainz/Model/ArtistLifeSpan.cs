namespace Brainz.Model
{
    using Newtonsoft.Json;

    public class ArtistLifeSpan
    {
        [JsonProperty("begin", NullValueHandling = NullValueHandling.Ignore)]
        public string Begin { get; set; }

        [JsonProperty("ended")]
        public string Ended { get; set; }

        [JsonProperty("end", NullValueHandling = NullValueHandling.Ignore)]
        public string End { get; set; }
    }
}
