namespace Brainz.Model
{
    using Newtonsoft.Json;

    public class Tag
    {
        [JsonProperty("count")]
        public long Count { get; set; }

        [JsonProperty("name")]
        public string Name { get; set; }
    }
}
