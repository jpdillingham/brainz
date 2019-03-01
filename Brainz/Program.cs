using Newtonsoft.Json;
using System;
using System.Linq;
using System.Net.Http;
using System.Threading.Tasks;
using Utility.CommandLine;
using System.Collections.Generic;
using Brainz.Responses;
using Brainz.Model;

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
        //private static Func<Guid, string> ReleaseGroupRequest => (mbid) => $"{ApiRoot}/artist/{mbid}?inc=release-groups&fmt=json";
        private static Func<Guid, int, int, string> ReleaseGroupRequest => (mbid, offset, limit) => $"{ApiRoot}/release-group?artist={mbid}&offset={offset}&limit={limit}&fmt=json";


        //private static Func<Guid, string> ReleaseRequest => (mbid) => $"{ApiRoot}/release-group/{mbid}?inc=releases+media&fmt=json";
        //https://musicbrainz.org/ws/2/release?release-group=04066384-affd-32ce-9b6f-ac4a97c07671&offset=0&limit=25&fmt=json&inc=media
        private static Func<Guid, int, int, string> ReleaseRequest => (mbid, offset, limit) => $"{ApiRoot}/release?release-group={mbid}&offset={offset}&limit={limit}&inc=media&fmt=json";
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

            Console.WriteLine();
            Console.WriteLine($"Fetching release group matches for artist '{bestArtist.Name}', album '{Album}'...");
            Console.WriteLine();

            List<ReleaseGroup> releaseGroups = new List<ReleaseGroup>();
            ReleaseGroupResponse releaseGroupResponse = null;

            do
            {
                req = ReleaseGroupRequest(bestArtist.Id, releaseGroups.Count, 100);
                var releaseGroupJson = await Http.GetStringAsync(req).ConfigureAwait(false);
                releaseGroupResponse = JsonConvert.DeserializeObject<ReleaseGroupResponse>(releaseGroupJson);

                releaseGroups.AddRange(releaseGroupResponse.ReleaseGroups);
            } while (releaseGroups.Count < releaseGroupResponse.ReleaseGroupCount);

            releaseGroups = releaseGroups.OrderByDescending(r => r.Title.SimilarityCaseInsensitive(Album)).ToList();

            var bestReleaseGroup = releaseGroups.FirstOrDefault();

            foreach (var releaseGroup in releaseGroups)
            {
                Console.WriteLine($"{(releaseGroup.Title == bestReleaseGroup.Title ? "-->" : "   ")} {(releaseGroup.Title.SimilarityCaseInsensitive(Album) * 100).ToString("F0").PadLeft(3)}%   {releaseGroup.Title}");
            }

            Console.WriteLine();
            Console.WriteLine($"Best release group match: {bestReleaseGroup.Title} (score: {(bestReleaseGroup.Title.SimilarityCaseInsensitive(Album) * 100).ToString("F0")}%)");


            //Console.WriteLine(req);
            Console.WriteLine();
            Console.WriteLine($"Fetching releases for release group '{bestReleaseGroup.Title}'...");
            Console.WriteLine();

            List<Release> releases = new List<Release>();
            ReleaseResponse releasesResponse = null;

            do
            {
                req = ReleaseRequest(bestReleaseGroup.Id, 0, 100);
                var releasesJson = await Http.GetStringAsync(req).ConfigureAwait(false);
                releasesResponse = JsonConvert.DeserializeObject<ReleaseResponse>(releasesJson);

                releases.AddRange(releasesResponse.Releases);
            } while (releases.Count < releasesResponse.ReleaseCount);
            
            releases = releases.OrderBy(r => r.Date.ToFuzzyDateTime()).ToList();

            var bestRelease = releases
                .Where(r => r.Status == "Official")
                .Where(r => r.Country == "US")
                .Where(r => string.IsNullOrEmpty(r.Disambiguation)).FirstOrDefault();

            foreach (var release in releases)
            {
                string trackStr = string.Join("+", release.Media.Select(m => m.TrackCount));
                string mediaStr = string.Join("+", release.Media.Select(m => m.Format));

                Console.WriteLine($"{(release.Id == bestRelease.Id ? "-->" : "   ")} {release.Title}\t{release.Country}\t{release.Date.PadRight(10)}\t{trackStr}\t{mediaStr}");
            }

            Console.WriteLine();
            Console.WriteLine($"Best release match: {bestRelease.Title}");

            req = RecordingRequest(bestRelease.Id);
            Console.WriteLine();
            Console.WriteLine($"Fetching recordings for release '{bestRelease.Title}' ({bestRelease.Id})");
            Console.WriteLine();

            var recordingJson = await Http.GetStringAsync(req).ConfigureAwait(false);
            var recordingResponse = JsonConvert.DeserializeObject<RecordingResponse>(recordingJson);

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
