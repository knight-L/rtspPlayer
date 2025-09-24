import { useState } from "react";
import WebRtcPlayer from "./WebRtcPlayer";

function App() {
  const [url, setUrl] = useState("");

  return (
    <div
      id="App"
      className="flex flex-col items-center justify-center gap-4 w-screen h-screen">
      <div className="w-2/3">
        <WebRtcPlayer url={url} />
      </div>
      <form
        id="myForm"
        onSubmit={(e) => {
          e.preventDefault();
          const formData = new FormData(e.currentTarget);
          const newUrl = formData.get("url") as string;
          setUrl(newUrl);
        }}>
        <input
          id="name"
          className="bg-[#222630] px-4 py-3 outline-none w-[280px] text-white rounded-lg border-2 transition-colors duration-100 border-solid focus:border-[#596A95] border-[#2B3040]"
          autoComplete="off"
          name="url"
          type="text"
          placeholder="rtsp地址"
          defaultValue={url}
        />
        <button
          type="submit"
          className="inline-block cursor-pointer items-center justify-center rounded-xl border-[1.58px] border-zinc-600 bg-zinc-800 px-5 py-3 font-medium text-slate-200 shadow-md ml-2 hover:bg-zinc-950 focus:bg-zinc-950">
          播放
        </button>
      </form>
    </div>
  );
}

export default App;
