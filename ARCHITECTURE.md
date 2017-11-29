IMDB Queue ---> Get HTML Document ---> Mux Fan Out ---> Append Movie IDs
                                            |
                                            |
                                             ---------> Extract Movie ------> Insert Movie


RT Queue

type App struct {
    IMDB *queue.Queue
    RT   *queue.Queue
    WG   sync.WaitGroup (?)
}

Mux stores map of host to channel slice (append and extract)

Append Movie IDs
Extract Movie
    Check doc.Url.Host for IMDB or RT
