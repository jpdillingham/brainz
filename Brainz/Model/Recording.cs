namespace brainz.Model
{
    using System;
    using Newtonsoft.Json;

    public partial class Recording
    {
        [JsonProperty("disambiguation")]
        public string Disambiguation { get; set; }

        [JsonProperty("id")]
        public Guid Id { get; set; }

        [JsonProperty("title")]
        public string Title { get; set; }

        [JsonProperty("video")]
        public bool Video { get; set; }

        [JsonProperty("length")]
        public long Length { get; set; }
    }
}