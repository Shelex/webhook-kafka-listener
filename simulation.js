import http from "k6/http";
import { check } from "k6";

export const options = {
  stages: [
    { duration: "5s", target: 1000 },
    { duration: "5s", target: 2000 },
    { duration: "10s", target: 4000 },
    { duration: "5s", target: 1000 },
    { duration: "5s", target: 0 },
  ],
};

export default function () {
  const url = "http://localhost:9000/webhooks";
  const payload = {
    updated: "2017-02-15T11:01:52.896Z",
    created: "2017-02-15T11:01:52.896Z",
    callback_id: "58a434bavclxkjsdf89d5a2b",
    owner_id: "582412sa87das87as68asd87wqe8be9d76",
    number: 9000,
    transaction_timestamp: "2017-02-15T11:01:52.722Z",
    id: "58a4352213123asdaseqwa355f46070",
  };

  const res = http.post(url, payload);
  check(res, { "status was 200": (r) => r.status == 200 });
}
