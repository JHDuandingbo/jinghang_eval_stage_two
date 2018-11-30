let startParam = {
  app: {
    userId: " guest", //用户ID，暂时不用, 必须
    secretId: " guest", //暂时不用,必须
    secretKey: "guest ", //暂时不用，必须
    securityMethod: "sha256" //暂时不用，必须
  },
  request: {
    coreType: "en.sent.score", //必须，评测类型, 待归纳整理,包括 句子评测、语义评测
    refText:
      "I would say yes and no. It depends on what kind of people you are. Surely, we live in a rapidly globalized world which is one free-flowing global labor market now. It is quite normal to work in a foreign country. But I think, only a few kinds of people can benefit from integration into the global world. What I mean is that you need to be really adaptable and skilled in order to overcome the stiff global competition, and not the other way around", //必须，评测参考文本
    precision: 0.5,
    phdet: 0,
    precision: 0.1,
    rank: 5,
    syldet: 0
  },
  /*
 [
  {
    "extra": "测评.part1.句子跟读",
    "weights": { "pron": 1, "stress": 1, "fluency": 1, "liaison": 0 },
    "requestKey": "evaluation.part1.sentReading"
  },
  {
    "extra": "测评.part2.情景问答",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanAccuracy": 1
    },
    "requestKey": "evaluation.part2.situationalQA"
  },
  {
    "extra": "测评.part2.答案跟读",
    "weights": { "pron": 1, "stress": 1, "fluency": 1, "liaison": 0 },
    "requestKey": "evaluation.part2.answerReading"
  },
  {
    "extra": "测评.part3.单词跟读",
    "weights": { "pron": 1, "stress": 0 },
    "requestKey": "evaluation.part3.wordReading"
  },
  {
    "extra": "测评.part3.单词例句跟读",
    "weights": { "pron": 1, "stress": 1, "fluency": 1, "liaison": 0 },
    "requestKey": "evaluation.part3.wordSentReading"
  },
  {
    "extra": "测评.part3.短语跟读",
    "weights": { "pron": 1, "stress": 1, "fluency": 1, "liaison": 0 },
    "requestKey": "evaluation.part3.phraseReading"
  },
  {
    "extra": "测评.part3.短语例句跟读",
    "weights": { "pron": 1, "stress": 1, "fluency": 1, "liaison": 0 },
    "requestKey": "evaluation.part3.phraseSentReading"
  },
  {
    "extra": "测评.part4.翻译",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanAccuracy": 1
    },
    "requestKey": "evaluation.part4.translation"
  },
  {
    "extra": "测评.part5.自由问答",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "vocabulary": 1,
      "grammar": 1,
      "relevantness": 0
    },
    "requestKey": "evaluation.part5.freeQA"
  },
  {
    "extra": "测评.part5.段落跟读",
    "requestKey": "evaluation.part5.paragraphReading"
  }
]


  */
  //requestKey编号规则:  测评或系统课.part名称.题目类型.题目编号
  //测评:evaluation  系统课:course
  // part名称: part1  part2  part3
  // 题目编号： 1 ,2,3,4,5
  requestKey: "course.part1.sentReading.1",
  sessionId: "xxxxxxxxxxxxxxxxxxxx",
  version: "0.1",
  platform: "android",
  compressed: 0, //0 :uncompressed , 1 :compressed
  action: "start", //必须
  userData: "fatcat user data",
  ts: 1537412219 //必须， 时间戳.
};

let rsp = {
  coreType: "en.sent.score",
  errId: 0,
  errMsg: null,
  result: {
    //兼容旧版
    scoreProFluency: "3.9",
    scoreProNoAccent: "3.9",
    scoreProStress: "3.9",

    //新版字段会改变
    sentence: "I'm looking forward to going abroad.",
    //同scoreProNoAccent
    pron: "3.9",
    //同scoreProStress
    stress: "3.9",
    //同scoreProFluency
    fluency: "3.9",
    //显示星星由overall决定
    overall: "3.9"
  },
  ts: "1543574796",
  userData: "",
  userId: "guest"
};
