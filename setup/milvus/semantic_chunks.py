from pymilvus import connections, Collection, CollectionSchema, FieldSchema, DataType, utility
import contextlib
from datasets import load_dataset
from langchain_experimental.text_splitter import SemanticChunker
from langchain_huggingface.embeddings import HuggingFaceEmbeddings
import time

port = 19530
dimensions = 384
max_url_length = 256
host = "192.168.1.160"
num_rows = 10

# Connect to Milvus
@contextlib.contextmanager
def milvus_connection(host, port):
    connections.connect("default", host=host, port=str(port))
    try:
        yield
    finally:
        connections.disconnect("default")

print("Connecting to Milvus")

with milvus_connection(host, 19530):

    collection_name = "chunks"
    id_col = "id"
    url_col = "url"
    embedding_col = "chunk_meaning"

    # Check if collection exists
    if utility.has_collection(collection_name):
        print("Pages collection already exists, dropping")
        Collection(collection_name).drop()

    # Define collection schema
    schema = CollectionSchema(
        fields=[
            FieldSchema(name=id_col, dtype=DataType.INT64, is_primary=True, auto_id=True),
            FieldSchema(name=url_col, dtype=DataType.VARCHAR, max_length=max_url_length),
            FieldSchema(name=embedding_col, dtype=DataType.FLOAT_VECTOR, dim=dimensions)
        ],
        description=f"{collection_name} schema"
    )

    # Create collection
    collection = Collection(name=collection_name, schema=schema, shards_num=2)
    print(f"Creating {collection_name} collection")

    if not collection.is_empty:
        print(f"Failed to create collection {collection_name}")
    else:
        print(f"Collection {collection_name} created successfully")

    #  -------------- SENTENCE EMBEDDINGS -------------------   #

    # Load the Wikipedia dataset from Hugging Face
    print("Loading dataset")
    ds = load_dataset("wikimedia/wikipedia", "20231101.en", split=f"train[:{num_rows}]")

    # Initialize SentenceTransformer embeddings
    print("Initializing SentenceTransformer embeddings")

    model_name = "sentence-transformers/all-MiniLM-L6-v2"
    model_kwargs = {'device': 'cpu'}
    encode_kwargs = {'normalize_embeddings': False}
    sentence_embeddings = HuggingFaceEmbeddings(
        model_name=model_name,
        model_kwargs=model_kwargs,
        encode_kwargs=encode_kwargs
    )

    # Split text using the SemanticChunker with SentenceTransformer embeddings
    print("Splitting text using SemanticChunker")
    text_splitter = SemanticChunker(
        sentence_embeddings,
        breakpoint_threshold_type="gradient",
        breakpoint_threshold_amount=85,
        buffer_size=3,
    )

    # Create documents with semantic chunking
    print("Creating documents")

    start_time = time.time()

    # iterate over rows in the dataset
    for row in ds:

        url = row.get('url', 'unknown')

        docs = text_splitter.create_documents([row['text']])

        # Extract chunked text contents
        doc_chunks = [doc.page_content for doc in docs]

        # Embed each chunk
        embeddings = sentence_embeddings.embed_documents(doc_chunks)

        data = [{"url": url, "chunk_meaning": embedding} for embedding in embeddings]
        
        # Insert data into Milvus collection
        collection.insert(data=data)

    print("Data insertion completed successfully.")
    print(f"Time taken: {time.time() - start_time} seconds")
