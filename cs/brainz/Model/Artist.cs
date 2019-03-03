namespace brainz.Model
{
    using System;
    using Newtonsoft.Json;

    public partial class Artist
    {
        [JsonProperty("id")]
        public Guid Id { get; set; }

        [JsonProperty("type", NullValueHandling = NullValueHandling.Ignore)]
        public string Type { get; set; }

        [JsonProperty("type-id", NullValueHandling = NullValueHandling.Ignore)]
        public Guid? TypeId { get; set; }

        [JsonProperty("score")]
        public long Score { get; set; }

        [JsonProperty("name")]
        public string Name { get; set; }

        [JsonProperty("sort-name")]
        public string SortName { get; set; }

        [JsonProperty("country", NullValueHandling = NullValueHandling.Ignore)]
        public string Country { get; set; }

        [JsonProperty("area", NullValueHandling = NullValueHandling.Ignore)]
        public Area Area { get; set; }

        [JsonProperty("begin-area", NullValueHandling = NullValueHandling.Ignore)]
        public Area BeginArea { get; set; }

        [JsonProperty("disambiguation", NullValueHandling = NullValueHandling.Ignore)]
        public string Disambiguation { get; set; }

        [JsonProperty("life-span")]
        public ArtistLifeSpan LifeSpan { get; set; }

        [JsonProperty("tags", NullValueHandling = NullValueHandling.Ignore)]
        public Tag[] Tags { get; set; }

        [JsonProperty("aliases", NullValueHandling = NullValueHandling.Ignore)]
        public Alias[] Aliases { get; set; }

        [JsonProperty("gender", NullValueHandling = NullValueHandling.Ignore)]
        public string Gender { get; set; }
    }
}
