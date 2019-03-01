namespace Brainz.Model
{
    using Newtonsoft.Json;
    using System;

    public class Area
    {
        [JsonProperty("name")]
        public string Name { get; set; }

        [JsonProperty("sort-name")]
        public string SortName { get; set; }

        [JsonProperty("iso-3166-1-codes")]
        public string[] Iso31661Codes { get; set; }

        [JsonProperty("id")]
        public Guid Id { get; set; }

        [JsonProperty("disambiguation")]
        public string Disambiguation { get; set; }

        [JsonProperty("type")]
        public string Type { get; set; }

        [JsonProperty("type-id")]
        public Guid TypeId { get; set; }

        [JsonProperty("life-span")]
        public AreaLifeSpan LifeSpan { get; set; }
    }
}
