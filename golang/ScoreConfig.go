package main


var configStr = `
[
  {
    "extra": "测评.part1.句子跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevantness": 0
    },
    "requestKey": "evaluation.part1.sentReading"
  },
  {
    "extra": "测评.part2.情景问答",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanAccuracy": 1,
      "vocabulary": 0,
      "grammar": 0,
      "relevantness": 0
    },
    "requestKey": "evaluation.part2.situationalQA"
  },
  {
    "extra": "测评.part2.答案跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevantness": 0
    },
    "requestKey": "evaluation.part2.answerReading"
  },
  {
    "extra": "测评.part3.单词跟读",
    "weights": {
      "pron": 1,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevantness": 0
    },
    "requestKey": "evaluation.part3.wordReading"
  },
  {
    "extra": "测评.part3.单词例句跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevantness": 0
    },
    "requestKey": "evaluation.part3.wordSentReading"
  },
  {
    "extra": "测评.part3.短语跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevantness": 0
    },
    "requestKey": "evaluation.part3.phraseReading"
  },
  {
    "extra": "测评.part3.短语例句跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevantness": 0
    },
    "requestKey": "evaluation.part3.phraseSentReading"
  },
  {
    "extra": "测评.part4.翻译",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanAccuracy": 1,
      "vocabulary": 0,
      "grammar": 0,
      "relevantness": 0
    },
    "requestKey": "evaluation.part4.translation"
  },
  {
    "extra": "测评.part5.自由问答",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanAccuracy": 0,
      "vocabulary": 1,
      "grammar": 1,
      "relevantness": 0
    },
    "requestKey": "evaluation.part5.freeQA"
  },
  {
    "extra": "测评.part5.段落跟读",
    "requestKey": "evaluation.part5.paragraphReading",
    "weights": {
      "pron": 1,
      "stress": 0,
      "fluency": 1,
      "liaison": 0,
      "semanAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevantness": 0
    }
  }
]
`



