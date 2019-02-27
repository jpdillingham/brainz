using Newtonsoft.Json;
using System;
using System.Linq;
using System.Net.Http;
using System.Threading.Tasks;
using Utility.CommandLine;
using Brainz.Artist;
using Brainz.ReleaseGroup;
using Brainz.Release;
using Brainz.Recording;

namespace Brainz
{
    class Program
    {
        [Argument('a', "artist", "The name of the artist to search for.")]
        private static string Artist { get; set; }

        [Argument('l', "album", "The name of the album to search for.")]
        private static string Album { get; set; }

        private static HttpClient Http { get; } = new HttpClient();

        private static string UserAgent = "Brainz/1.00 (https://github.com/jpdillingham/Brainz)";
        private static string ApiRoot = @"https://musicbrainz.org/ws/2";

        private static Func<string, string> ArtistRequest => (artist) => $"{ApiRoot}/artist/?query={artist}&fmt=json";
        private static Func<Guid, string> ReleaseGroupRequest => (mbid) => $"{ApiRoot}/artist/{mbid}?inc=release-groups&fmt=json";
        private static Func<Guid, string> ReleaseRequest => (mbid) => $"{ApiRoot}/release-group/{mbid}?inc=releases+media&fmt=json";
        private static Func<Guid, string> RecordingRequest => (mbid) => $"{ApiRoot}/release/{mbid}?inc=recordings&fmt=json";

        static int Main(string[] args)
        {
            Arguments.Populate();

            Http.DefaultRequestHeaders.UserAgent.ParseAdd(UserAgent);

            return Task.Run(() => MainAsync(args)).ConfigureAwait(false).GetAwaiter().GetResult();
        }

        private static async Task<int> MainAsync(string[] args)
        {
            var req = ArtistRequest(Artist);
            Console.WriteLine($"Fetching artist matches for '{Artist}'...");
            Console.WriteLine();

            var artistJson = await Http.GetStringAsync(req).ConfigureAwait(false);
            var artistResponse = JsonConvert.DeserializeObject<ArtistResponse>(artistJson);
            var artists = artistResponse.Artists.OrderByDescending(a => a.Score);

            var bestArtist = artists.FirstOrDefault();

            foreach (var artist in artistResponse.Artists)
            {
                Console.WriteLine($"{(artist.Name == bestArtist.Name ? "-->" : "   ")} {artist.Score.ToString().PadLeft(3)}%   {artist.Name}");
            }

            Console.WriteLine();
            Console.WriteLine($"Best artist match: {bestArtist.Name} (score: {bestArtist.Score}%)");

            req = ReleaseGroupRequest(bestArtist.Id);
            Console.WriteLine();
            Console.WriteLine($"Fetching release group matches for artist '{bestArtist.Name}', album '{Album}'...");
            Console.WriteLine();

            var releaseGroupJson = await Http.GetStringAsync(ReleaseGroupRequest(bestArtist.Id)).ConfigureAwait(false);
            var releaseGroupResponse = JsonConvert.DeserializeObject<ReleaseGroupResponse>(releaseGroupJson);
            var releaseGroups = releaseGroupResponse.ReleaseGroups.OrderByDescending(r => r.Title.SimilarityCaseInsensitive(Album));

            var bestReleaseGroup = releaseGroups.FirstOrDefault();

            foreach (var releaseGroup in releaseGroups)
            {
                Console.WriteLine($"{(releaseGroup.Title == bestReleaseGroup.Title ? "-->" : "   ")} {(releaseGroup.Title.SimilarityCaseInsensitive(Album) * 100).ToString("F0").PadLeft(3)}%   {releaseGroup.Title}");
            }

            Console.WriteLine();
            Console.WriteLine($"Best release group match: {bestReleaseGroup.Title} (score: {(bestReleaseGroup.Title.SimilarityCaseInsensitive(Album) * 100).ToString("F0")}%)");

            req = ReleaseRequest(bestReleaseGroup.Id);
            Console.WriteLine();
            Console.WriteLine($"Fetching releases for release group '{bestReleaseGroup.Title}'...");
            Console.WriteLine();

            var releasesJson = await Http.GetStringAsync(req).ConfigureAwait(false);
            var releasesResponse = JsonConvert.DeserializeObject<ReleaseResponse>(releasesJson);
            var releases = releasesResponse.Releases.OrderBy(r => r.Date.ToFuzzyDateTime());

            Console.WriteLine(releasesJson);

            var bestRelease = releases
                .Where(r => r.Status == "Official")
                .Where(r => r.Country == "US")
                .Where(r => string.IsNullOrEmpty(r.Disambiguation)).FirstOrDefault();

            foreach (var release in releases)
            {
                Console.WriteLine($"{release.Title}\t{release.Country}\t{release.Date}\t{release");
            }

            Console.WriteLine($"Best release match: {bestRelease.Title}");

            req = RecordingRequest(bestRelease.Id);
            Console.WriteLine($"Fetching {req}...");
            var recordingJson = await Http.GetStringAsync(req).ConfigureAwait(false);
            var recordingResponse = JsonConvert.DeserializeObject<RecordingResponse>(recordingJson);
            Console.WriteLine($"{recordingResponse.Media[0].Tracks.Count()} recordings returned.");

            foreach (var media in recordingResponse.Media)
            {
                foreach (var recording in media.Tracks.OrderBy(r => r.Position))
                {
                    Console.WriteLine($"{media.Position}{recording.Position.ToString("D2")} - {recording.Title}");
                }
            }

            return 0;
        }
    }
}
