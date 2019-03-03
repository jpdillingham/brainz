using System;
using System.Collections.Generic;
using System.Reflection;
using System.Text;

namespace brainz
{
    public static class Extensions
    {
        public static DateTime ToFuzzyDateTime(this string s)
        {
            if (s.Length == 2)
            {
                if (int.Parse(s) < 30)
                {
                    return DateTime.Parse($"1-1-20{s}");
                }

                return DateTime.Parse($"1-1-19{s}");
            }

            if (s.Length == 4)
            {
                return DateTime.Parse($"1-1-{s}");
            }

            else return DateTime.Parse(s);
        }

        public static double Similarity(this string s, string t)
        {
            return (1.0 - ((double)s.LevenshteinDistance(t) / (double)Math.Max(s.Length, t.Length)));
        }

        public static double SimilarityCaseInsensitive(this string s, string t)
        {
            return (1.0 - ((double)s.LevenshteinDistanceCaseInsensitive(t) / (double)Math.Max(s.Length, t.Length)));
        }

        public static int LevenshteinDistanceCaseInsensitive(this string s, string t)
        {
            return s.ToLower().LevenshteinDistance(t.ToLower());
        }

        public static int LevenshteinDistance(this string s, string t)
        {
            int n = s.Length;
            int m = t.Length;
            int[,] d = new int[n + 1, m + 1];

            if (n == 0)
            {
                return m;
            }

            if (m == 0)
            {
                return n;
            }

            for (int i = 0; i <= n; d[i, 0] = i++)
            {
            }

            for (int j = 0; j <= m; d[0, j] = j++)
            {
            }

            for (int i = 1; i <= n; i++)
            {
                for (int j = 1; j <= m; j++)
                {
                    int cost = (t[j - 1] == s[i - 1]) ? 0 : 1;

                    d[i, j] = Math.Min(
                        Math.Min(d[i - 1, j] + 1, d[i, j - 1] + 1),
                        d[i - 1, j - 1] + cost);
                }
            }

            return d[n, m];
        }
    }
}
