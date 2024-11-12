from datasets import load_dataset
from langchain_experimental.text_splitter import SemanticChunker
from langchain_huggingface.embeddings import HuggingFaceEmbeddings
import time

num_rows = 10

# Load the Wikipedia dataset from Hugging Face
print("Loading dataset")
ds = load_dataset("wikimedia/wikipedia", "20231101.en", split=f"train[:{num_rows}]")
text_data = ds['text']  # Using a subset for demonstration; you can adjust this

print("Text data:")
print(text_data)

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

docs = text_splitter.create_documents(text_data)

end_time = time.time()

print(f"-------------- Created {len(docs)} Documents--------------")

for doc in docs:
    print("\n")
    print("\n")
    print(doc.page_content)
    print("-----------------------------")

print(f"Time taken: {end_time - start_time} seconds. Average time per document: {(end_time - start_time) / num_rows} seconds.")