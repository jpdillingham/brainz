namespace brainz.Model
{
    using Newtonsoft.Json;

    public class Alias
    {
        [JsonProperty("sort-name")]
        public string SortName { get; set; }

        [JsonProperty("name")]
        public string Name { get; set; }

        [JsonProperty("locale")]
        public object Locale { get; set; }

        [JsonProperty("type")]
        public object Type { get; set; }

        [JsonProperty("primary")]
        public object Primary { get; set; }

        [JsonProperty("begin-date")]
        public object BeginDate { get; set; }

        [JsonProperty("end-date")]
        public object EndDate { get; set; }
    }
}
