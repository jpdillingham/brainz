namespace Brainz.Model
{
    using System;
    using Newtonsoft.Json;

    public class Track
    {
        [JsonProperty("position")]
        public long Position { get; set; }

        [JsonProperty("id")]
        public Guid Id { get; set; }

        [JsonProperty("length")]
        public long Length { get; set; }

        [JsonProperty("title")]
        public string Title { get; set; }

        [JsonProperty("number")]
        public string Number { get; set; }

        [JsonProperty("recording")]
        public Recording Recording { get; set; }
    }
}
