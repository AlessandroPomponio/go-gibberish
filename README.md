# Go gibberish

[![Actions Status](https://github.com/AlessandroPomponio/go-gibberish/workflows/Go/badge.svg)](https://github.com/AlessandroPomponio/go-gibberish/actions)
[![Build Status](https://travis-ci.org/AlessandroPomponio/go-gibberish.svg?branch=master)](https://travis-ci.org/AlessandroPomponio/go-gibberish)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlessandroPomponio/go-gibberish)](https://goreportcard.com/report/github.com/AlessandroPomponio/go-gibberish)
[![GoDoc](https://godoc.org/github.com/AlessandroPomponio/go-gibberish?status.svg)](https://godoc.org/github.com/AlessandroPomponio/go-gibberish)

This program is a go-powered version of <https://github.com/rrenaud/Gibberish-Detector>.

It uses a training file to build a model, which is then used to check whether a string is likely to be gibberish or not.

## How it works

With rrenaud's words, the creator of the original Python algorithm:

> It uses a 2 character markov chain.
>
> The markov chain first 'trains' or 'studies' a few MB of English text, recording how often characters appear next to each other. Eg, given the text "Rob likes hacking" it sees Ro, ob, o[space], [space]l, ... It just counts these pairs. After it has finished reading through the training data, it normalizes the counts. Then each character has a probability distribution of 27 followup character (26 letters + space) following the given initial.
>
>So then given a string, it measures the probability of generating that string according to the summary by just multiplying out the probabilities of the adjacent pairs of characters in that string. EG, for that "Rob likes hacking" string, it would compute prob['r']['o'] * prob['o']['b'] * prob['b'][' '] ... This probability then measures the amount of 'surprise' assigned to this string according the data the model observed when training. If there is funny business with the input string, it will pass through some pairs with very low counts in the training phase, and hence have low probability/high surprise.
>
>I then look at the amount of surprise per character for a few known good strings, and a few known bad strings, and pick a threshold between the most surprising good string and the least surprising bad string. Then I use that threshold whenever to classify any new piece of text.
>
>Peter Norvig, the director of Research at Google, has this nice talk about "The unreasonable effectiveness of data" here, <http://www.youtube.com/watch?v=9vR8Vddf7-s>. This insight is really not to try to do something complicated, just write a small program that utilizes a bunch of data and you can do cool things.

## How to use it

Run the training for the model by calling the function `training.TrainModel` and then use `gibberish.IsGibberish` to detect whether a string is gibberish or not.
In case you decide to us

```go

    var (
        performTraining bool
    )

    func main() {

        flag.BoolVar(&performTraining, "train", false, "train")
        flag.Parse()
        
        if performTraining {
            err := training.TrainModel(consts.AcceptedCharacters, "big.txt", "good.txt", "bad.txt", "knowledge.json")
            if err != nil {
                log.Fatal(err)
            }
            
            return
        }
        
        reader := bufio.NewReader(os.Stdin)
        data, err := persistence.LoadKnowledgeBase("knowledge.json")
        if err != nil {
        	log.Fatal(err)
        }
        
        for {
        
        	fmt.Print("Insert something to check: ")
        	input, _ := reader.ReadString('\n')
        	input = strings.TrimSpace(input)
        	isGibberish := gibberish.IsGibberish(input, data)
        	fmt.Println(fmt.Sprintf("Input: %s: is gibberish? %v\n", input, isGibberish))
        
        }

    }
```

## Credits

Thanks once again to [rrenaud](https://github.com/rrenaud) for the original algorithm.

A huge thank you goes to [domef](https://github.com/domef) as well, for helping me translate the algorithm.
