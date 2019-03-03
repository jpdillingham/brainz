namespace brainz.Model
{
    using Newtonsoft.Json;
    using System;

    public class Media
    {
        [JsonProperty("format-id")]
        public Guid? FormatId { get; set; }

        [JsonProperty("track-count")]
        public long TrackCount { get; set; }

        [JsonProperty("title")]
        public string Title { get; set; }

        [JsonProperty("position")]
        public long Position { get; set; }

        [JsonProperty("format")]
        public string Format { get; set; }

        [JsonProperty("tracks")]
        public Track[] Tracks { get; set; }

        [JsonProperty("track-offset")]
        public long TrackOffset { get; set; }
    }
}
