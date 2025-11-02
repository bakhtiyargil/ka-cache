import http from "k6/http";
import {check, sleep} from "k6";
import {SharedArray} from "k6/data";
import {Trend, Counter} from "k6/metrics";

const config = JSON.parse(open("./config.json"));
const {baseUrl, endpoints, keys, stages, thresholds, defaultTtl} = config;

export let options = {stages, thresholds};

const putTrend = new Trend("PUT_duration");
const getTrend = new Trend("GET_duration");
const cacheHits = new Counter("cache_hits");
const cacheMisses = new Counter("cache_misses");

const keyPool = new SharedArray("keys", () =>
    Array.from({length: keys}, (_, i) => `key-${i}`)
);

export default function () {
    const key = keyPool[Math.floor(Math.random() * keyPool.length)];
    const op = Math.random();

    if (op < 0.7) {
        //this is for "http_req_failed" case
        http.setResponseCallback(http.expectedStatuses(200, 404));

        const res = http.get(`${baseUrl}${endpoints.get}/${key}`);
        getTrend.add(res.timings.duration);
        check(res, {"GET status 200 or 404": (r) => [200, 404].includes(r.status)});

        if (res.status === 200) cacheHits.add(1);
        else cacheMisses.add(1);
    } else {
        const payload = JSON.stringify({
            key,
            value: `value-${__VU}-${__ITER}`,
            ttl: defaultTtl
        });

        const res = http.put(`${baseUrl}${endpoints.put}`, payload, {
            headers: {"Content-Type": "application/json"},
        });

        putTrend.add(res.timings.duration);
        check(res, {"PUT status 200": (r) => [200].includes(r.status)});
    }

    sleep(0.001);
}
