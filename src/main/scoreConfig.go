package main

var ScoreConfigStr = `
{
  "evaluation.part1.sentReading": {
    "extra": "测评.part1.句子跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part1.sentReading"
  },
  "evaluation.part1.wordReading": {
    "extra": "测评.part1.单词跟读",
    "weights": {
      "pron": 1,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part1.wordReading"
  },
  "evaluation.part2.translation": {
    "extra": "测评.part2.翻译",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 1,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part2.translation"
  },

  "evaluation.part2.situationalQA": {
    "extra": "测评.part2.情景问答",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 1,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part2.situationalQA"
  },
  "evaluation.part2.answerReading": {
    "extra": "测评.part2.答案跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part2.answerReading"
  },
  "evaluation.part3.translation": {
    "extra": "测评.part3.翻译",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 1,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part3.translation"
  },
  "evaluation.part3.wordReading": {
    "extra": "测评.part3.单词跟读",
    "weights": {
      "pron": 1,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part3.wordReading"
  },
  "evaluation.part3.wordSentReading": {
    "extra": "测评.part3.单词例句跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part3.wordSentReading"
  },
  "evaluation.part3.phraseReading": {
    "extra": "测评.part3.短语跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part3.phraseReading"
  },
  "evaluation.part3.phraseSentReading": {
    "extra": "测评.part3.短语例句跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part3.phraseSentReading"
  },
"evaluation.part4.sentReading": {
    "extra": "测评.part4.句子跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part4.sentReading"
  },

  "evaluation.part4.translation": {
    "extra": "测评.part4.翻译",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 1,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "evaluation.part4.translation"
  },
  "evaluation.part5.freeQA": {
    "extra": "测评.part5.自由问答",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 1,
      "grammar": 1,
      "relevancy": 0
    },
    "requestKey": "evaluation.part5.freeQA"
  },
  "evaluation.part5.paragraphReading": {
    "extra": "测评.part5.段落跟读",
    "requestKey": "evaluation.part5.paragraphReading",
    "weights": {
      "pron": 1,
      "stress": 0,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancrelevancy": 0
    }
  },
  "course.part1.sentReading": {
    "extra": "系统课.part1.句子跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "course.part1.sentReading"
  },
  "course.part2.situationalQA": {
    "extra": "系统课.part2.情景问答",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 1,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "course.part2.situationalQA"
  },
  "course.part2.answerReading": {
    "extra": "系统课.part2.答案跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "course.part2.answerReading"
  },
  "course.part3.wordReading": {
    "extra": "系统课.part3.单词跟读",
    "weights": {
      "pron": 1,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "course.part3.wordReading"
  },
  "course.part3.wordSentReading": {
    "extra": "系统课.part3.单词例句跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "course.part3.wordSentReading"
  },
  "course.part3.phraseReading": {
    "extra": "系统课.part3.短语跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "course.part3.phraseReading"
  },
  "course.part3.phraseSentReading": {
    "extra": "系统课.part3.短语例句跟读",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "course.part3.phraseSentReading"
  },
  "course.part4.translation": {
    "extra": "系统课.part4.翻译",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 1,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "course.part4.translation"
  },
  "course.part5.freeQA": {
    "extra": "系统课.part5.自由问答",
    "weights": {
      "pron": 0,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 1,
      "grammar": 1,
      "relevancy": 0
    },
    "requestKey": "course.part5.freeQA"
  },
  "course.part5.paragraphReading": {
    "extra": "系统课.part5.段落跟读",
    "requestKey": "course.part5.paragraphReading",
    "weights": {
      "pron": 1,
      "stress": 0,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    }
  },
	"ifun.italk.dub": {
    "extra": "ifun.italk.配音挑战",
    "requestKey": "ifun.italk.dub",
    "weights": {
      "pron": 1,
      "stress": 1,
      "fluency": 1,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    }
  },
"ieltsword": {
    "extra": "雅思词汇",
    "weights": {
      "pron": 1,
      "stress": 0,
      "fluency": 0,
      "liaison": 0,
      "semanticAccuracy": 0,
      "vocabulary": 0,
      "grammar": 0,
      "relevancy": 0
    },
    "requestKey": "ieltsword"
  }


}
`

/*
func main(){
	var   obj map[string]interface{}


    println(obj["heello"] == nil)



}
*/
