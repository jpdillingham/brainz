namespace brainz.Responses
{
    using System;
    using brainz.Model;
    using Newtonsoft.Json;

    public class ArtistResponse
    {
        [JsonProperty("created")]
        public DateTimeOffset Created { get; set; }

        [JsonProperty("count")]
        public long Count { get; set; }

        [JsonProperty("offset")]
        public long Offset { get; set; }

        [JsonProperty("artists")]
        public Artist[] Artists { get; set; }
    }
}
