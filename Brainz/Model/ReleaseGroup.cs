namespace Brainz.Model
{
    using System;
    using Newtonsoft.Json;

    public partial class ReleaseGroup
    {
        [JsonProperty("secondary-type-ids")]
        public Guid[] SecondaryTypeIds { get; set; }

        [JsonProperty("disambiguation")]
        public string Disambiguation { get; set; }

        [JsonProperty("first-release-date")]
        public string FirstReleaseDate { get; set; }

        [JsonProperty("primary-type-id")]
        public Guid? PrimaryTypeId { get; set; }

        [JsonProperty("primary-type")]
        public string PrimaryType { get; set; }

        [JsonProperty("id")]
        public Guid Id { get; set; }

        [JsonProperty("title")]
        public string Title { get; set; }

        [JsonProperty("secondary-types")]
        public string[] SecondaryTypes { get; set; }
    }
}