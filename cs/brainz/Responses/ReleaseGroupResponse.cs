namespace brainz.Responses
{
    using brainz.Model;
    using Newtonsoft.Json;

    public partial class ReleaseGroupResponse
    {
        [JsonProperty("release-groups")]
        public ReleaseGroup[] ReleaseGroups { get; set; }

        [JsonProperty("release-group-count")]
        public long ReleaseGroupCount { get; set; }

        [JsonProperty("release-group-offset")]
        public long ReleaseGroupOffset { get; set; }
    }
}
